# MySQL 索引学习、总结

以下实验、练习均在`MySQL8.0`上进行，8.0 是MySQL的一个重大版本，虽然用的人较少。

数据库版本：

```sh
Server version: 8.0.27
```

## 索引分类

MySQL索引包括普通索引、唯一索引、全文索引、单列索引、多列索引和空间索引等。

- 从`功能`上划分，索引可分为：普通索引、唯一索引、主键索引和全文索引4种；
- 从`物理实现`上划分，可分为：聚簇索引和非聚簇索引两大类；
- 从`组织字段`上划分，可分为：单列索引和多列索引两部分。

当然索引之间是存在重合的，如主键索引也是一种唯一索引，索引键是不可重复的；不同的  
存储引擎支持的索引也不大相同：

- InnoDB：支持 bree，full-text索引(5.7后支持)，不支持 hash 索引；
- MyISAM：支持 bree，full-text索引，不支持 hash 索引；
- Memory：支持 btree、hash 等索引，不支持 full-text 索引;
- And so on...

### 索引创建

索引既可以在创建表时顺带创建，也可以后面单独创建。

如果在创建表时就创建索引，其基本语法如下：

```
CREATE TABLE table_name [col_name data_type]
[UNIQUE | FULLTEXT | SPATIAL] [INDEX | KEY] [index_name] (col_name [length]) [ASC |
DESC]
```

- UNIQUE(唯一)、FULLTEXT(全文)、SPATIAL(空间) 作为索引修饰词，表明索引的类型；
- INDEX，KEY 都可用来表示索引，同义词；
- index_name 索引名称，如果不指定，默认以字段名称作为索引名；
- col_name 需要创建索引的字段名称，可以为多个；
- length 索引长度，只有字符串类型(非定长)需要指定索引长度；
- ASC、DESC用来指定索引排序方式。

例子：

```mysql
CREATE TABLE test_index1
(
    id   INT(11)  NOT NULL AUTO_INCREMENT,
    name CHAR(30) NOT NULL,
    age  INT(11)  NOT NULL,
    info VARCHAR(255),
    geo  GEOMETRY NOT NULL,

    KEY name_idx (name (20)),
    INDEX multi_idx (id, name, age),
    FULLTEXT INDEX futxt_idx_info (info),
    SPATIAL INDEX spa_idx_geo (geo),
    PRIMARY KEY (id)
);
```

这个例子中新建了各种类型的索引：

- id 为自增主键索引；
- name_idx 为长度为20的姓名索引；
- futxt_idx_info 为全文搜索索引；
- spa_idx_geo 为空间索引。

已存在的表，可使用`alter`和`create index`来新建索引，语法如下：

```sql
 ALTER TABLE table_name ADD [UNIQUE | FULLTEXT | SPATIAL] [INDEX | KEY]
[index_name] (col_name[length],...) [ASC | DESC]

CREATE
[UNIQUE | FULLTEXT | SPATIAL] INDEX index_name
ON table_name (col_name[length],...) [ASC | DESC]
```

字段说明与上文一致。

### 索引删除

每一个索引都是要付出代价的，插入一条数据记录，伴随着大量的索引创建工作，对于非必要的  
索引，可以使用`alter`和`drop index`来删掉，语法如下：

```sql
ALTER TABLE table_name DROP INDEX index_name;

DROP INDEX index_name ON table_name;
```

> 提示：
> 1. 如果列在索引中，当该列被删除时，索引中的列也会删除；如果组成索引的所有列都被删除了，
     > 那么整个索引也会被删除；
> 2. 降序索引是8.0才开始支持的；

### 索引隐藏

MySQL在8.0后开吃支持索引隐藏；索引一旦被设置为隐藏，那么优化器将不再选择该索引，即使该  
索引被 force index 强制使用，优化器也不会选择该索引；将索引设置为隐藏后，系统在一段时间  
内正常运行，那么证明该索引是无用的，考虑此时再来彻底删除该索引。

同样地，可以通过`alter`来设置一个索引是否隐藏：

```sql
#切换成隐藏索引
ALTER TABLE tablename ALTER INDEX index_name INVISIBLE;

#切换成非隐藏索引
ALTER TABLE tablename ALTER INDEX index_name VISIBLE;
```

当然，事无绝对，如果隐藏的索引在业务突然被需要，更改索引可见行未免麻烦，因此  
MySQL提供了配置`use_invisible_indexes`来设置隐藏索引是否可用，如果可用，那么  
优化器仍有可能选中被隐藏的索引，如下：

```mysql
set session optimizer_switch = "use_invisible_indexes=on";
```

## 索引基本原则

1. 字段有唯一性限制；
2. 频繁作为 where 查询、更新、删除的条件；
3. 频繁作为 group by、order by 使用的列；
4. 需要 distinct 的字段需建索引；
5. 如果字符串字段较长，使用字符串前缀作为索引，   
   可以使用`select count(distinct left(address,10)) / count(*)`来计算前缀的覆盖率，   
   选择合适的前缀即可；
6. 根据索引的最左前缀原则，使用最频繁的字段放在最左边；
7. 多个字段都需索引的情况下，优先考虑联合索引；
8. 有大量重复值的字段不适合建索引，如 `gender` 字段；
9. 无序的字段不适合建索引，如UUID，MD5等；
10. 索引不是银弹，需注重业务架构的良好设计。

## 索引优化

### 索引失效

1. 按照最左前缀原则，如果左边的值未确定，那么其它的值也无法使用相关的索引；
2. 表达式计算、函数调用、类型转换会导致索引失效；
3. is null可以使用索引，而 is not null无法使用索引；
4. like以通配符`%`开头，索引会失效；
5. OR的左、右两侧存在所索引列，索引会失效；

### join 优化

1. 保证被驱动表的 join 字段已有索引；
2. left join 时，选择小表作为驱动表，大表作为被驱动表；
3. 减少子查询使用，改用 join 优化(自连接)；

### order by 优化

1. sql 中，可以在 where 子句和 order by 子句中使用索引，目的是避免全表扫描和 file sort；
2. 尽量使用索引完成 order by 排序，如果 where 和 order by 后面是相同的列就使用单索引列; 如果不同就使用联合索引。
3. order by 无法使用索引时，需要对 file sort 方式进行调优：
    1. 对 where 和 order by 中的字段新建一个联合索引；
    2. 尝试提高 sort_buffer_size；
    3. 尝试提高 max_length_for_sort_data；

### group by 优化

1. group by 也可使用索引，先排序后分组；
2. 当无法使用索引时，也会 file sort，尝试提高 sort_buffer_size、max_length_for_sort_data；
3. where 的效率高于 having，优先将筛选条件写在 where 中。

### 索引覆盖

当能通过读取索引就可以得到想要的数据，那就不需要回表；一个索引包含了满足查询结果的数据就叫做覆盖索引。

简单地说，索引中的字段已经覆盖了 select 要查询的字段。

优点：

1. 避免回表后二次查询；
2. 减少IO次数。

### 前缀索引

1. 选择区分度高的索引长度，减少数据重复；
2. 索引长度能节省一定的空间；
3. 前缀索引无用使用索引覆盖，因为索引数据只是一部分，需要回表拿到全部。

### 其它

- 索引下推
- 唯一索引、普通索引
- so on...