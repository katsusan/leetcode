参考链接：
  mysql的事务隔离级别和锁：https://tech.meituan.com/2014/08/20/innodb-lock.html


1. Basic
  - 数据库level
  ```
  CREATE DATABASE `db` character set utf8;
  DROP DATABASE `db`;
  
  ```

  GROUP BY
  ------------
  - The GROUP BY Clause is used to group rows with same values .
  - The GROUP BY Clause is used together with the SQL SELECT statement.
  - The SELECT statement used in the GROUP BY clause can only be used contain column names, aggregate functions, constants and expressions.
  - The HAVING clause is used to restrict the results returned by the GROUP BY clause.


2. explain相关
    from： https://segmentfault.com/a/1190000008131735

   各个字段含义：
    - id: SELECT 查询的标识符. 每个 SELECT 都会自动分配一个唯一的标识符.

    - select_type: SELECT 查询的类型.
        + SIMPLE, 表示此查询不包含 UNION 查询或子查询

        + PRIMARY, 表示此查询是最外层的查询

        + UNION, 表示此查询是 UNION 的第二或随后的查询

        + DEPENDENT UNION, UNION 中的第二个或后面的查询语句, 取决于外面的查询

        + UNION RESULT, UNION 的结果

        + SUBQUERY, 子查询中的第一个 SELECT

        + DEPENDENT SUBQUERY: 子查询中的第一个 SELECT, 取决于外面的查询. 即子查询依赖于外层查询的结果.

    - table: 查询的是哪个表

    - partitions: 匹配的分区

    - type: join 类型
        ALL < index < range ~ index_merge < ref < eq_ref < const < system

        + system: 表中只有一条数据. 这个类型是特殊的 const 类型.

        + const: 针对主键或唯一索引的等值查询扫描, 最多只返回一行数据. const 查询速度非常快, 因为它仅仅读取一次即可.
          例如下面的这个查询, 它使用了主键索引, 因此 type 就是 const 类型的.

        + eq_ref: 此类型通常出现在多表的 join 查询, 表示对于前表的每一个结果, 都只能匹配到后表的一行结果. 并且查询的比较操作通常是 =, 查询效率较高. 

        + ref: 此类型通常出现在多表的 join 查询, 针对于非唯一或非主键索引, 或者是使用了 最左前缀 规则索引的查询. 

        + range: 表示使用索引范围查询, 通过索引字段范围获取表中部分数据记录. 这个类型通常出现在 =, <>, >, >=, <, <=, IS NULL, <=>, BETWEEN, IN()   操作中.
          当 type 是 range 时, 那么 EXPLAIN 输出的 ref 字段为 NULL, 并且 key_len 字段是此次查询中使用到的索引的最长的那个.

        + index: 表示全索引扫描(full index scan), 和 ALL 类型类似, 只不过 ALL 类型是全表扫描, 而 index 类型则仅仅扫描所有的索引, 而不扫描数据.
          index 类型通常出现在: 所要查询的数据直接在索引树中就可以获取到, 而不需要扫描数据. 当是这种情况时, Extra 字段 会显示 Using index.
        
        + ALL: 表示全表扫描, 这个类型的查询是性能最差的查询之一. 通常来说, 我们的查询不应该出现 ALL 类型的查询, 因为这样的查询在数据量大的情况下, 对   数据库的性能是巨大的灾难. 如一个查询是 ALL 类型查询, 那么一般来说可以对相应的字段添加索引来避免.

    - possible_keys: 此次查询中可能选用的索引

    - key: 此次查询中确切使用到的索引.

    - key_len: 表示查询优化器使用了索引的字节数. 这个字段可以评估组合索引是否完全被使用, 或只有最左部分字段被使用到.

    - ref: 哪个字段或常数与 key 一起被使用

    - rows: 显示此查询一共扫描了多少行. 这个是一个估计值. （rows是核心指标，绝大部分rows小的语句执行一定很快，所以优化语句基本上都是在优化rows）

    - filtered: 表示此查询条件所过滤的数据的百分比

    - extra: 额外的信息
        + Using filesort： 当 Extra 中有 Using filesort 时, 表示 MySQL 需额外的排序操作, 不能通过索引顺序达到排序效果. 一般有 Using filesort, 
          都建议优化去掉, 因为这样的查询 CPU 资源消耗大.

        + Using index： "覆盖索引扫描", 表示查询在索引树中就可查找所需数据, 不用扫描表数据文件, 往往说明性能不错

        + Using temporary：查询有使用临时表, 一般出现于排序, 分组和多表 join 的情况, 查询效率不高, 建议优化.







3. 索引设计
  添加索引命令：
    - ALTER TABLE `ground` ADD PRIMARY KEY (`uid`);   //添加主键索引
    - ALTER TABLE `ground` ADD UNIQUE (`uid`);  //添加唯一索引
    - ALTER TABLE `ground` ADD INDEX index_create_time (`create_time`); //添加单一索引
    - ALTER TABLE `ground` ADD FULLTEXT (`body`); //添加全文索引
    - ALTER TABLE `ground` ADD INDEX index_id_name (`id`, `name`);  //添加联合索引

