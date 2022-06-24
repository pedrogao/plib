# LSM tree

## 核心设计点

1. LSM 日志结构存储引擎，核心点在于顺序写日志，这样就能保证超快的写速度；
2. 每个日志文件都有长度限制(threshold)，超过这个限制则新建另一个日志文件写；
3. 在内存中保存了索引、数据方便快速查询，如果仍查不到则去搜索日志文件；
4. 如果一个 key 被写了多次，那么就会有很多重复的行，因此需要合并他们(compact)；
5. 引入布隆过滤器来加快文件数据查询，不存在的 key 直接返回，避免读文件；
6. memtable 保存了 k-v 信息，方便快速查询，sparse index 保存了部分 k-file-offset信息，方便在文件中搜索；
7. 稀疏索引查询的时候找到比 key 小（对于二叉树的左孩子节点）的文件和索引，然后遍历文件，就能提升效率；
8. WAL 日志用来恢复 memtable；
9.

## references

- [LSM-Tree](https://github.com/chrislessard/LSM-Tree)