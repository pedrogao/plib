#!/usr/bin/env python3

class LockManager:
    """
    lock manager for transactions
    """

    def __init__(self):
        self.locks = []

    def add(self, transaction, record_id):
        if not self.exists(transaction, record_id):
            self.locks.append([transaction, record_id])

    def exists(self, transaction, record_id):
        return any(lock[0] is transaction and lock[1] == record_id for lock in self.locks)


class Table:
    """
    table, set of records
    """

    def __init__(self):
        self.next_xid = 0
        self.active_xids = set()  # active transaction ids
        self.records = []  # record
        self.locks = LockManager()  # lock manager

    def new_transaction(self, transaction_type):
        self.next_xid += 1
        self.active_xids.add(self.next_xid)  # add to active transactions
        return transaction_type(self, self.next_xid)


class RollbackException(Exception):
    pass


class Transaction:
    def __init__(self, table, xid):
        self.table = table  # table
        self.xid = xid
        self.rollback_actions = []  # rollback action

    def add_record(self, id, name):
        record = {
            'id': id,
            'name': name,
            'created_xid': self.xid,  # created transaction id
            'expired_xid': 0  # expired transaction id
        }
        self.rollback_actions.append(["delete", len(self.table.records)])  # reverse
        self.table.records.append(record)

    def delete_record(self, rid):
        for i, record in enumerate(self.table.records):
            if self.record_is_visible(record) and record['id'] == rid:
                if self.record_is_locked(record):
                    raise RollbackException("Row locked by another transaction.")
                else:
                    record['expired_xid'] = self.xid
                    self.rollback_actions.append(["add", i])

    def update_record(self, rid, name):
        self.delete_record(rid)
        self.add_record(rid, name)

    def fetch_record(self, rid):
        for record in self.table.records:
            if self.record_is_visible(record) and record['id'] is rid:
                return record

        return None

    def count_records(self, min_id, max_id):
        count = 0
        for record in self.table.records:
            if self.record_is_visible(record) and \
                    min_id <= record['id'] <= max_id:
                count += 1

        return count

    def fetch_all_records(self):
        visible_records = []
        for record in self.table.records:
            if self.record_is_visible(record):
                visible_records.append(record)

        return visible_records

    def fetch(self, expr):
        visible_records = []
        for record in self.table.records:
            if self.record_is_visible(record) and expr(record):
                visible_records.append(record)

        return visible_records

    def commit(self):
        self.table.active_xids.discard(self.xid)

    def rollback(self):
        for action in reversed(self.rollback_actions):
            if action[0] == 'add':
                self.table.records[action[1]]['expired_xid'] = 0
            elif action[0] == 'delete':
                self.table.records[action[1]]['expired_xid'] = self.xid

        self.table.active_xids.discard(self.xid)

    def record_is_visible(self, record):
        raise NotImplementedError()

    def record_is_locked(self, record):
        raise NotImplementedError()


class ReadUncommittedTransaction(Transaction):
    """
    read un-committed
    """

    def record_is_locked(self, record):
        # 如果 expired_xid 不为 0，即已经被删除，则被锁住了(已经删除了还看个啥)
        return record['expired_xid'] != 0

    def record_is_visible(self, record):
        # expired_xid 为0，表示记录还在 MVCC 中，可见
        return record['expired_xid'] == 0


class ReadCommittedTransaction(Transaction):
    """
    read committed
    """

    def record_is_locked(self, record):
        # 如果 expired_xid 不为 0，即已经被删除，且不在当前活跃的事务中
        # 即使记录被删除了，如果仍在活跃事务列表中，那么记录被锁住了
        return record['expired_xid'] != 0 and \
               record['expired_xid'] in self.table.active_xids

    def record_is_visible(self, record):
        # The record was created in active transaction that is not our
        # own.
        if record['created_xid'] in self.table.active_xids and \
                record['created_xid'] != self.xid:
            return False

        # The record is expired or and no transaction holds it that is
        # our own.
        if record['expired_xid'] != 0 and \
                (record['expired_xid'] not in self.table.active_xids or
                 record['expired_xid'] == self.xid):
            return False

        return True


class RepeatableReadTransaction(ReadCommittedTransaction):
    """
    read repeated
    """

    def record_is_locked(self, record):
        return ReadCommittedTransaction.record_is_locked(self, record) or \
               self.table.locks.exists(self, record['id'])

    def record_is_visible(self, record):
        is_visible = ReadCommittedTransaction.record_is_visible(self, record)

        if is_visible:
            self.table.locks.add(self, record['id'])

        return is_visible


class SerializableTransaction(RepeatableReadTransaction):
    """
    serialize
    """

    def __init__(self, table, xid):
        Transaction.__init__(self, table, xid)
        self.existing_xids = self.table.active_xids.copy()

    def record_is_visible(self, record):
        is_visible = ReadCommittedTransaction.record_is_visible(self, record) \
                     and record['created_xid'] <= self.xid \
                     and record['created_xid'] in self.existing_xids

        if is_visible:
            self.table.locks.add(self, record['id'])

        return is_visible


class TransactionTest:
    """
    transaction test baseline
    """

    def __init__(self, transaction_type):
        self.table = Table()
        # 新增两条记录
        client = self.table.new_transaction(ReadCommittedTransaction)
        client.add_record(id=1, name="Joe")
        client.add_record(id=3, name="Jill")
        client.commit()

        self.client1 = self.table.new_transaction(transaction_type)
        self.client2 = self.table.new_transaction(transaction_type)

    def run_test(self):
        try:
            return self.run()
        except RollbackException:
            return False

    def result(self):
        if self.run_test():
            return '✔'
        return '✘'

    def run(self):
        raise NotImplementedError()


class DirtyRead(TransactionTest):
    def run(self):
        # 事务1，读取 record1
        result1 = self.client1.fetch_record(rid=1)
        # 事务2，更新 record1
        self.client2.update_record(rid=1, name="Joe 2")
        # 事务1，再次读取 record1
        result2 = self.client1.fetch_record(rid=1)
        # 如果 result1 和 result2 不相等，则发生了脏读
        return result1 != result2


class NonRepeatableRead(TransactionTest):
    def run(self):
        # 事务1，读取 record1
        result1 = self.client1.fetch_record(rid=1)
        # 事务2，更新 record1
        self.client2.update_record(rid=1, name="Joe 2")
        # 事务2，提交
        self.client2.commit()
        # 事务1，再次读取 record1
        result2 = self.client1.fetch_record(rid=1)
        # 如果 result1 和 result2 不相等，则表示发生了不可重复读
        return result1 != result2


class PhantomRead(TransactionTest):
    def run(self):
        # 事务1，读取 record1,  record2
        result1 = len(self.client1.fetch(lambda r: 1 <= r['id'] <= 3))
        # 事务2，新增 record3
        self.client2.add_record(id=2, name="John")
        # 事务2，提交
        self.client2.commit()
        # 事务1，再次读取 record1,  record2
        result2 = self.client1.count_records(min_id=1, max_id=3)
        # 如果 result1 和 result2 不相等，则表示发生了幻读
        return result1 != result2


if __name__ == '__main__':
    isolation_modes = [
        ['read uncommitted', ReadUncommittedTransaction],
        ['read committed  ', ReadCommittedTransaction],
        ['repeatable read ', RepeatableReadTransaction],
        ['serializable    ', SerializableTransaction]
    ]

    possible_errors = [DirtyRead, NonRepeatableRead, PhantomRead]

    print('                  Dirty Repeat Phantom')
    for isolation_mode in isolation_modes:
        results = [possible_error(isolation_mode[1]).result() for possible_error in possible_errors]
        print(isolation_mode[0] + "    " + results[0] + "      " + results[1] + "      " + results[2])
