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

4. InnoDB
  show engine InnoDB status;

5. MySQL排序实现
  若语句里有order by或者group by这样的操作，代表需要对输出结果进行排序。
  MySQL大概有两类排序方式：

  - 使用range、ref、index之类。explain输出上述代表对索引列的读取方式，这样读出来不需要额外的排序操作。

  - 使用filesort算法，简而言之就是将一组记录按照快速排序算法放入内存然后再归并，filesort所能使用的内存
      由@@sort_buffer_size决定。如果排序数据大于@@sort_buffer_size则会创建一个临时表存储使得性能下降。

    filesort的数据通常来自一个表，如果来自多个表则将多个表的内容放入临时表，再用filesort对它进行操作。
    其工作模式有两种：
		+ 直接模式：要求排序的数据已经被完全读取，排完序的数据就是我们需要的。
		+ 间接模式：排序时数据以<sort_key, rowid>的形式存在，排完后根据rowid找到数据再读取。
			(比如被排序的记录包含Blob等变长字段时会使用间接模式)
    

6. MySQL网络协议
	- 数据交互流程
		+ 握手阶段(客户端开始连接时)：
			* 服务器->客户端：握手初始化包
			* 客户端->服务器：客户端认证包
			* 服务器->客户端：OK、Error包
		+ 命令包(客户端对服务器任一请求)：
			* 客户端->服务器：命令包
			* 服务器->客户端：OK、Error、结果集包
	- 网络包格式
		+ 字符串:
			* 以NULL结尾
			* 带长度标识的字符串(Length Coded String)
				- 第一字节0~250 => 代表后续数据最多250字节
				- 第一字节251	=> 列值为NULL，仅用于行数据包
				- 第一字节252	=> 后续2个字节的值代表字符串有多少字节数据
				- 第一字节253	=> 后续3个字节的值代表字符串有多少字节数据
				- 第一字节254	=> 后续8个字节的值代表字符串有多少字节数据
		+ 网络包头部
			4个字节的包头，前3个字节代表MySQL数据包头部之后的数据长度，也就是限制了最大2^24B=16MB，
			后面1个字节代表数据包序列号，每个命令开始后它都会被重置为0。
	- 客户端包格式
		+ 认证包
			* 4字节：客户端标志
			* 4字节：包最大长度
			* 1字节：客户端字符集ID
			* 23字节：填充字符0x00
			* N字节：用户名，以NULL结尾的字符串
			* 1+N字节：密码加密字段(sramble_buff)，带长度标识的字符串
			* N字节：数据库名
		+ 命令包
			1字节命令+N字节命令参数
			命令如：SLEEP，QUIT，SHUTDOWN，STMT_PREPARE，STMT_EXECUTE等等

			show global status like 'Com%'可以获得服务器接收到的各种命令统计。

			例：use test = 0x02(use) 0x74 0x65 0x73 0x74(test)
	- 服务端包格式
		+ 



