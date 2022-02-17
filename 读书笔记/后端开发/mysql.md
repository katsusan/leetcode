参考链接：
  mysql的事务隔离级别和锁：https://tech.meituan.com/2014/08/20/innodb-lock.html
  阿里数据库blog：http://mysql.taobao.org/monthly

# 0. Pre
  install from source
  - groupadd mysql
  - useradd -r -g mysql -s /bin/false mysql
  - tar zxvf mysql-XXXX.tar.gz
  - cd mysql-XXXX
  - mkdir bld
  - cd bld
  - make clean && rm CmakeCache.txt   //optional, if 
  - cmake .. -DWITH_DEBUG=1 -DWIN_DEBUG_NO_INLINE=1 -D-DWITH_INNODB_EXTRA_DEBUG=1 -DWITH_BOOST=/usr/local/boost_1_59_0
  - make && make install DESTDIR="/usr/local/mysql"
  - /usr/local/mysql/scripts/mysql_install_db --user=mysql    //create my.cnf

  configure
  - set password = password('123456');   //first time when you login into mysql
  - grant all on db1.tb1 to 'user1'@'host1';

# 1. Basic
  - 数据库level
  ```
  CREATE DATABASE `db` character set utf8;
  DROP DATABASE `db`;
  
  ```

  - grant
    - CREATE USER 'jeffrey'@'localhost' IDENTIFIED BY 'password';
    - GRANT ALL ON db1.* TO 'jeffrey'@'localhost';
    - GRANT SELECT ON db2.invoice TO 'jeffrey'@'localhost';
    - ALTER USER 'jeffrey'@'localhost' WITH MAX_QUERIES_PER_HOUR 90;
    


  GROUP BY
  ------------
  - The GROUP BY Clause is used to group rows with same values .
  - The GROUP BY Clause is used together with the SQL SELECT statement.
  - The SELECT statement used in the GROUP BY clause can only be used contain column names, aggregate functions, constants and expressions.
  - The HAVING clause is used to restrict the results returned by the GROUP BY clause.


# 2. explain相关
    from： https://segmentfault.com/a/1190000008131735
    better ref: https://dev.mysql.com/doc/refman/5.7/en/explain-output.html

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

        + range: 表示使用索引范围查询, 通过索引字段范围获取表中部分数据记录. 这个类型通常出现在 =, <>, >, >=, <, <=, IS NULL, <=>, BETWEEN, IN()操作中.
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



# 3. 索引设计
  索引操作：
    - ALTER TABLE `ground` ADD PRIMARY KEY (`uid`);   //添加主键索引
    - ALTER TABLE `ground` ADD UNIQUE (`uid`);  //添加唯一索引
    - ALTER TABLE `ground` ADD INDEX index_create_time (`create_time`); //添加单一索引
    - ALTER TABLE `ground` ADD FULLTEXT (`body`); //添加全文索引
    - ALTER TABLE `ground` ADD INDEX index_id_name (`id`, `name`);  //添加联合索引
    - ALTER TABLE `ground` DROP INDEX `index_id_name`;  // 删除索引index_id_name

# 4. InnoDB
  show engine InnoDB status;

# 5. MySQL排序实现
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
    

# 6. MySQL网络协议
	- 数据交互流程
		+ 握手阶段(客户端开始连接时)：
      * 客户端->服务器：请求连接
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
		+ 握手初始化包
      * 1字节：协议版本号，由include/mysql_version.h中的PROTOCOL_VERSION定义
      * n字节：服务器信息，n=strlen(server_version)+1, MYSQL_SERVER_VERSION中保存着版本信息，比如"5.0.77-community-nt MySQL Community Edition(GPL)"
      * 4字节：线程号，MySQL为此连接分配的线程号
      * 8字节：密码验证1
      * 1字节：0x00填充
      * 2字节：以位图形式表达出服务端可接受的连接选项，在include/mysql_com.h中以CLIENT_开头列出
              比如#define CLIENT_LONG_PASSWORD	1 /* new more secure passwords */
                  #define CLIENT_FOUND_ROWS	2	/* Found instead of affected rows */
                  #define CLIENT_LONG_FLAG	4	/* Get all column flags */
              注：头文件里CLIENT_开头的宏最大到了1<<31,因此这里怀疑应该是4字节
      * 1字节：服务器字符集
      * 2字节：服务器状态标志，在mysql_com.h中以SERVER_STATUS_开头，
              比如#define SERVER_STATUS_IN_TRANS     1
                  #define SERVER_STATUS_AUTOCOMMIT   2	/* Server in auto_commit mode */
                  #define SERVER_MORE_RESULTS_EXISTS 8    /* Multi query - next query exists */
      * 13字节：0x00填充
      * 13字节：密码验证2
    + 结果包
      开头1字节定义了包类型
        * 0x00：OK包
        * 0xff：ERROR包
        * 0xfe：EOF包
        * 1-250：结果集包，Select * from table1
        * 1-250：属性包， Select 1+1
        * 1-250：行数据包
      OK包：
        MySQL成功执行一个命令后，将回复一个OK包，通常是对以下客户端命令的回复：
          * COM_PING
          * COM_QUERY //这里的query指的是更广泛意义的查询，包括insert、update、delete等
          * COM_REFRESH
          * COM_REGISTER_SLAVE
        包括：1字节(field_count)+1~9字节(影响行数)+1~9字节(插入ID)+2字节(服务器状态)+2字节(警告数量)+N字节(消息)
      ERROR包：
        MySQL服务器处理命令出错，或者认证信息有问题则会返回ERROR包，包括：
          1字节(0xff)+2字节(错误号)+1字节(SQL状态标识符'#')+5字节(SQL状态)+N字节(消息)
      结果集包：
        宏观上由一系列包组成：
          * 结果集包头部->1~9字节(field_count)+1~9字节(附属字段)
          * 属性包1，属性包2，...，属性包N
          * EOF包
          * 行数据包1，行数据包2，...，行数据包N
          * EOF包
        查询结果如果是n列，m行，则最后结果集共有m+n+3个包。

# 7. MySQL查询解析与优化器
  MySQL解析器由两部分组成：
    - 词法分析(Lexical Analysis/Scanner)
        扫描输入的字符流，根据构词规则识别单词符号。可以用lex或GNU开源的flex完成。
    
    - 语法分析(Syntax Analysis/Parser)
        将识别的单词符号组合成各类语法短语，比如程序，语句，表达式等。可以用Bison等处理。
        见sql/sql_yacc.yy
    
    $lex是词法分析程序的自动生成工具，接受输入为描述构词规则的正规式，构建有穷自动机进而生成一个词法分析程序$
    $YACC(Yet Another Compiler Compiler)是Unix上一个用来生成编译器的编译器，采用LALR(1)语法分析方法$

    //TODO: 梳理解析器流程
  
  查询优化器：
    关系数据库管理系统早期的两个基本原型：INGRES和SYSTEM-R分别采用两种不同的查询优化技术：
      - 查询分解算法：Query Decompos，一种启发式的"贪婪"算法。
        + 如果查询有三个或更多的变量or关系，先用启发规则将查询划分为两个更小的子查询，
          若子查询包含三个以上则继续分解，否则用元组替代(Tuple Substitution)来处理。
          使用单变量查询处理器来创建用于扫描单个关系的路径。
      - 穷尽搜索算法：
        + 对于查询内出现的每个关系进行扫描，找出扫描的所有可能方法，对每个方法创建一条路径。
          通过使用这些路径考虑二元连接的所有计划。然后单个关系再形成三元连接，重复上述过程直到
          所有关系都今夕连接并形成路径。估计所有路径的执行代价然后选择最便宜的路径来创建执行计划。
          缺点在于随着连接数目变化复杂度呈指数增加，优点在于考虑了所有路径可以得到最优解。
      MySQL采用的是穷尽搜索算法，一般分为以下4个步骤：
        - 将查询转化为内部表示形式，通常是语法树
        - 根据一定的等价规则将语法树转化为标准形式
        - 选择底层算法，对于语法树中的每个操作根据存取路径、数据的存储发布等信息选择具体的执行算法
        - 生成查询计划，查询计划由一些内部操作组成，需要对其进行代价评估选择最小成本的执行方案
   
  查询的具体实现在handle_query函数，一个simple_query的调用链大概如下：
  handle_connection
    do_command
      dispatch_command
        mysql_parse
          mysql_execute_command
            execute_sqlcom_select
              handle_query
                select->prepare(thd)  //准备阶段，解析所有的表/列信息以及所有语句，改造语法树为之后的优化做准备(位于sql_resolver.cc)
                lock_tables(thd, lex->query_tables, lex->table_count, 0)
                query_cache.store_query(thd, lex->query_tables)
                select->optimize(thd) //优化的核心逻辑
                select->join->exec()  //执行查询,内部调用了do_select(位于sql_executor.cc)

# 8. MySQL安全管理
  - 默认用户管理
    - drop user 'user1'@'host1';
    - set password for 'user2'@'host2' = password('user2pwd');
    - mysqladmin -u root password 'newpwd'  // 修改root的密码为newpwd
  - 权限表的导入导出
    - select * from mysql.user into outfile '/tmp/user.txt';  // 将user表的所有内容导入到文件/tmp/user.txt
    - system more /tmp/user.txt;    // 查看user.txt内容
    - load data infile '/etc/password' into table test.os_passwd; // 将/etc/passwd导入到test.os_passwd表
    也就是只要登录用户有File_priv，可以随意对有对应权限的表导入导出。
    
    - $ mysqld_safe --skip-grant-tables   // 绕过权限表启动，root密码丢失时可以用这个重置，当然本身也有风险
      可以在编译mysqld时指定--disable-grant-options来禁用这个选项。

# 9. 经典存储引擎
MySQL的存储元文件位于datadir指定的选项下，为每个DB创建一个独立的文件夹，即$datadir/$dbname，
各个存储引擎创建的文件各不相同，但大多会用frm文件存储表结构，包括下面要说的MyISAM和InnoDB。

# 9.1 MyISAM
MyISAM类型的表会创建三个文件，tb.frm/tb.MYD/tb.MYI，分别存储表结构/数据/索引，具体如下：
  - MYD文件：数据文件，与InnoDB不同。MyISAM并不采用分页方式存取数据，也就不会在各行数据之间看到填充行。
    格式分为动态格式、固定格式、压缩格式。压缩格式只能由myisampack工具来创建。
    - 固定格式：记录头+数据1+记录头+数据2+...
    - 动态格式：当表中每行数据变化较大时，动态格式存储可以节省很多空间，特别是包含blob和varchar等变长类型。
    // TODO: 具体格式有空再探究

# 9.2 InnoDB
# 9.2.1 特性
  - 行级锁
  - SELECT下提供MVCC(多版本并发控制)
  - 双写入：进行表空间数据写操作时会将数据写两次(日志只写一次)
  - 插入缓存：减少了对非主键索引插入数据时造成的磁盘随机读写
  - 适应式哈希索引：InnoDB会监视对表定义的索引的检索，如果可以从建立哈希索引中收益则会自动建立它

# 9.2.2 架构
  ----------------------------------
  MySQL Server      应用
  ↓↑                 ↓↑
  Handle API        嵌入式InnoDB API
  ----------------------------------
              事务
  ----------------------------------
            游标/行
  Mini事务    B树       锁
              页
  ----------------------------------
              缓存
  ----------------------------------
        文件空间管理器
  ----------------------------------
              IO
  ----------------------------------            

  - 对于事务层，在InnoDB中所有的行为都发生在事务中，如果auto-commit被启用，则执行的每一个语句都是一个单独的事务
  - 锁功能层：完成锁功能和事务管理(比如回滚、提交等操作)。InnoDB还特意用一个表来来跟踪各种锁的情况。

# 9.2.3 InnoDB文件格式
InnoDB的表空间可以由多个文件或裸分区组成。如果innodb_file_per_table为OFF，则所有InnoDB的表共享一个表空间。
innodb_data_home_dir为物理存储路径。
表空间至少包含一下几个部分：
  - 内部数据词典
  - Undo日志
  - 插入缓存
  - 双写入缓存
  - MySQL复制功能相关信息

InnoDB中表空间-段-集合与行之间的关系：
  表空间：叶子节点段、非叶子节点段、回滚段
    ↓
  段：扩展块1、扩展块2、
    ↓
  扩展块：页1、页2、
    ↓
  页：行1、行2、
    ↓
  行：事务id、回滚指针、字段指针、字段1、字段2、字段3、

# 9.2.4 InnoDB记录结构
物理记录的组成：
  - 字段偏移量：1或2字节，记录的总长度小于127时为1字节，额外字段里会指出具体是1字节还是2字节
  - 额外字段：6字节
  - 记录内容: ----

字段偏移量(field start offsets)：
字段偏移量是一个目录列表，每个目录是一个相对于记录起点的偏移量，这些目录反向存储，即第一个字段的偏移量是存储在列表的最后一个目录。
比如1Byte字段+2Byte字段+4Byte字段，各自的偏移量分别为1，1+2，1+2+4，但因为是反向存储故字段偏移量里存的是0x07,0x03,0x01。

  1字节偏移量：
  - 第一个比特位：0-代表字段非NULL，1-代表NULL
  - 剩余7位：实际偏移值，0~127
  2字节偏移量：
  - 第一个比特位：0-代表字段非NULL，1-代表NULL
  - 第二个比特位：0-代表字段存储在同一页，1-代表存储在不同页，当包含大字段对象时可能会出现这一情况
  - 剩余14位：实际偏移值，0~16383

额外字节(extra bytes)：
定长6字节，包括:
  - 2bit：保留
  - 1bit：删除标志
  - 1bit：最小长度记录，1代表这一行为表中的最小行
  - 4bit：包含的数量
  - 13bit：记录在索引页中的位置
  - 10bit：列数，1~1023
  - 1bit：偏移长度标志
  - 16bit：指针，指向同页的下一个记录

当创建一个InnoDB表时，会额外创建三个系统列：行ID、事务ID、回滚指针。

# 9.2.5 InnoDB页结构
InnoDB将所有记录存储在数据库页中，一般非压缩的页大小为16KB。

一个InnoDB页包含7个部分：
  - 文件头(Fil Header):
  - 页头(Page Header)
  - 最小及最大虚记录(Infimum+Supremum Records)
  - 用户记录(User Records): 用户所有的插入数据，该区域不会按照B树的Key进行排序(避免频繁的数据移动)，插入时是插入到现有行的后面或者使用空闲空间
      但由于要按照Key值排序因此每个记录都包含一个next指针指向按key顺序的下一条记录。
  - 自由堆(Free Space)：
  - 页目录(Page Directory)
  - 文件尾(Fil Trailer)：

# 10. MySQL日志功能及实现分析
# 10.1 错误日志
默认情况下错误日志会记录到数据目录，名为<hostname>.err。
在命令行或配置文件里以log-error选项指定，可以用show variables like 'log_error'来查看。
// TODO(次要)：梳理错误日志初始化以及写逻辑。
// 初始化：mysqld_main() -> init_error_log(),简单初始化 -> init_server_components(),打开日志文件、初始化handler等
// 写：sql_print_error/sql_print_warning/sql_print_information -> error_log_print -> log_error -> vprint_msg_to_log -> print_buffer_to_file

# 10.2 普通日志
general query log记录了连接、断开连接、以及每次执行的sql语句。
默认情况下是未启用状态，可以用以下方式启用：
  - at startup：general_log=1或无->启用 | 0代表禁用；general_log_file=filename，指定日志存储位置
  - at runtime：set global general_log = 'ON'/1 |'OFF'/0, set global general_log_file = '/path/general.log'
// TODO(次要)：梳理普通日志初始化及写逻辑
// 初始化：mysqld_main() -> query_logger.init()，初始化全局变量query_logger(负责general/slow log) -> init_server_components()
//           -> query_logger.set_handlers(log_output_options)/reopen_log_file(QUERY_LOG_GENERAL)
// 写：query_logger.general_log_print -> general_log_write -> (Log_to_csv_event_handler/Log_to_file_event_handler)::log_general

# 10.3 慢查询日志
超过long_query_time(默认10秒)的语句会被记录在慢查询日志，格式如下：
```sql
# Time: 2020-12-13T18:06:49.198960+08:00
# User@Host: eleusr[eleusr] @  [192.168.1.33]  Id:    19
# Query_time: 10.000782  Lock_time: 0.000000 Rows_sent: 1  Rows_examined: 0
SET timestamp=1607854009;
select sleep(10);
```
默认慢查询日志不会记录，可以用以下方式启用：
  - at startup：slow_query_log=0(禁用)|1或无(启用); slow_query_log_file='/path',指定慢查询日志存储位置
  - at runtime: set global slow_query_log = ON/1 | 'OFF'/0；set global slow_query_log_file = '/path'
// TODO(次要)：梳理初始化及写逻辑
// 初始化：mysqld_main()  -> init_server_components() -> query_logger.set_handlers(log_output_options)/reopen_log_file(QUERY_LOG_SLOW)
// 写：log_slow_statement -> log_slow_do -> query_logger.slow_log_write -> slow_log_handler_list.log_slow
        -> (Log_to_csv_event_handler/Log_to_file_event_handler)::log_slow

# 10.4 二进制日志
二进制日志记录所有更新的SQL语句，包含可能更新数据的语句，当然SELECT和show语句是不会记录在内的。
它主要有两个功能，数据恢复和数据复制。默认情况下MySQL不会记录binlog，可以设置log-bin=[basename]
来启用。生成的binlog文件为basename.00000X格式，每次mysql启动都会令X加1。可以用mysqlbinlog命令读取。
启用后会产生basename.index和basename.00000X文件，前者存放既存binlog的路径，后者为实际的binlog文件。
注：设置log-bin参数时server-id也要一并设置。

// 摘取了一条插入语句`insert into account values (12, 'Trump', 8000)`的binlog内容，
```
#201213 19:18:17 server id 1  end_log_pos 293 CRC32 0x352cae79  Query   thread_id=2     exec_time=0     error_code=0
SET TIMESTAMP=1607858297/*!*/;
SET @@session.pseudo_thread_id=2/*!*/;
SET @@session.foreign_key_checks=1, @@session.sql_auto_is_null=0, @@session.unique_checks=1, @@session.autocommit=1/*!*/;
SET @@session.sql_mode=1436549152/*!*/;
SET @@session.auto_increment_increment=1, @@session.auto_increment_offset=1/*!*/;
/*!\C utf8 *//*!*/;
SET @@session.character_set_client=33,@@session.collation_connection=33,@@session.collation_server=45/*!*/;
SET @@session.lc_time_names=0/*!*/;
SET @@session.collation_database=DEFAULT/*!*/;
BEGIN
/*!*/;
# at 293
#201213 19:18:17 server id 1  end_log_pos 349 CRC32 0xff8fd1f2  Table_map: `ground`.`account` mapped to number 112
# at 349
#201213 19:18:17 server id 1  end_log_pos 400 CRC32 0xf2567dfd  Write_rows: table id 112 flags: STMT_END_F
BINLOG '
efjVXxMBAAAAOAAAAF0BAAAAAHAAAAAAAAEABmdyb3VuZAAHYWNjb3VudAADAw8DAvwDBvLRj/8=
efjVXx4BAAAAMwAAAJABAAAAAHAAAAAAAAEAAgAD//gMAAAABQBUcnVtcEAfAAD9fVby
'/*!*/;
```
// TODO: binlog的初始化及写逻辑
// 初始化：mysqld_main()  -> init_server_components() -> mysql_bin_log.open_index_file(opt_binlog_index_name, ln, TRUE)，打开index文件
//            -> tc_log= &mysql_bin_log; tc_log->open(opt_bin_log ? opt_bin_logname : opt_tc_log_file)
// 写：

# 11. 子系统
# 11.1 复制功能子系统(Replication)

# 11.2 错误消息子系统





