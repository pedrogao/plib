#!/usr/bin/env python3

next_xid = 1
active_xids = set()
records = []


def new_transaction():
    global next_xid
    next_xid += 1
    active_xids.add(next_xid)
    return Transaction(next_xid)


class Transaction:
    def __init__(self, xid):
        self.xid = xid
        self.rollback_actions = []

    def add_record(self, record):
        record['created_xid'] = self.xid
        record['expired_xid'] = 0
        self.rollback_actions.append(["delete", len(records)])
        records.append(record)

    def delete_record(self, id):
        for i, record in enumerate(records):
            if self.record_is_visible(record) and record['id'] == id:
                if self.row_is_locked(record):
                    raise "Row locked by another transaction."
                else:
                    record['expired_xid'] = self.xid
                    self.rollback_actions.append(["add", i])  # 该事务可能会回滚删除

    def update_record(self, id, name):
        self.delete_record(id)  # 先删除
        self.add_record({"id": id, "name": name})  # 后新增，可能会回滚，记录更改

    def row_is_locked(self, record):
        # 记录：未过期，且在活跃事务里面
        return record['expired_xid'] != 0 and record['expired_xid'] in active_xids

    def record_is_visible(self, record):
        # 总结一下，什么时候记录是可见的，按照版本来说：
        # - expired_xid 小于当前 xid 的记录可不见；
        # - created_xid 大于当前 xid 的记录不可见；

        # The record was created in active transaction that is not our
        # own.
        # 事务活跃，但不是自己；如果 created_xid 还在活跃列表里面，证明记录被创建了，但未被提交，所以事务只有自己可见
        if record['created_xid'] in active_xids and record['created_xid'] != self.xid:
            return False

        # The record is expired or and no transaction holds it that is
        # our own.
        # 已过期，不在活跃列表 ｜ 自己是自己删的
        # 如果事务已过期，但 expired_xid 仍在活跃列表中，那么记录仍是可见的；
        # 或者事务已过期，但 expired_xid 是当前事务，即记录是自己删除的，那么记录是不可见的；
        if record['expired_xid'] != 0 and \
                (record['expired_xid'] not in active_xids or
                 record['expired_xid'] == self.xid):
            return False

        return True

    def fetch_all(self):
        # 获取所有可见的记录
        global records
        visible_records = []
        for record in records:
            if self.record_is_visible(record):
                visible_records.append(record)
        return visible_records

    def commit(self):
        # 提交记录
        active_xids.discard(self.xid)

    def rollback(self):
        # 回滚
        for action in reversed(self.rollback_actions):
            if action[0] == 'add':
                records[action[1]]['expired_xid'] = 0  # 设置为0，即未删除
            elif action[0] == 'delete':
                records[action[1]]['expired_xid'] = self.xid

        active_xids.discard(self.xid)


if __name__ == '__main__':
    client1 = new_transaction()
    client2 = new_transaction()
    client1.add_record({"id": 1, "name": "Bob"})
    client1.add_record({"id": 2, "name": "John"})
    client1.delete_record(id=1)
    client1.update_record(id=2, name="Tom")
    print(client1.fetch_all())
    print(client2.fetch_all())

    client1.rollback()
    print(client2.fetch_all())
