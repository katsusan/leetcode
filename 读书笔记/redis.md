redis:
    任务队列：  RPUSH+BLPOP / LPUSH+BRPOP
    优先队列：  BLPOP queue1 queue2 ...

redis设计与实现：
    哈希表：
        扩展时空间变为大于等于ht[0].used * 2的最小2的幂数，
        收缩时空间变为大于等于ht[0].used的最小2的幂数。

        负载因子 load_factor = ht[0].used / ht[0].size  //已保存结点数量/哈希表大小
        - 未执行BGSAVE/BGREWRITEAOF时，load_factor >= 1 则扩展哈希表
        - 执行BGSAVE/BGREWRITEAOF时，load_factor >= 5 则扩展哈希表  (鉴于子进程的copy on write策略，提高要求避免不必要的内存操作)
        - load_factor <= 0.1则收缩哈希表

        渐进式哈希：ht[0]到ht[1]的迁移并非一次性完成，而是每次对字典有增删查改时将ht[0]的rehashidx索引上的键值对复制到ht[1]，
                并将rehashidx加1。最终复制完rehashidx变成-1.
                期间查找先查ht[0]再查ht[1]，新增键值对一律加入ht[1].
    
        结构(dict.h)：
            typedef struct dict {   //采用hashtable时的robj对象里的ptr所指
                dictType *type;
                void *privdata;
                dictht ht[2];
                long rehashidx; /* rehashing not in progress if rehashidx == -1 */
                unsigned long iterators; /* number of iterators currently running */
            } dict;

            typedef struct dictht {
                dictEntry **table;
                unsigned long size;
                unsigned long sizemask;
                unsigned long used;
            } dictht;

            typedef struct dictEntry {
                void *key;
                union {
                    void *val;
                    uint64_t u64;
                    int64_t s64;
                    double d;
                } v;
                struct dictEntry *next;
            } dictEntry;


    整数集合(intset):
        集合只包含整数且元素数量不多时就会使用intset作为集合的底层实现。
        结构：  
            typedef struct intset {
                //编码方式，sizeof(int16_t->2 /int32_t->4 /int64_t->8 )
                uint32_t encoding;
                //集合长度
                uint32_t length;
                //保存元素的底层数组
                int8_t contents[];
            } intset;
    
    压缩列表(ziplist)： //哈希键和列表键的可用选择之一
        结构：
            zlbytes: uint32_t, 整个压缩列表占用的字节数
            zltail: uint32_t, 表尾结点距起始地址有多少个字节
            zllen: uinte16_t, 节点数量(当小于65535时，为真实节点数量。当等于65535时需要遍历才能得出真实节点数量)
            entry1
            ...     //zlentry, 各个节点
            entryN
            zlend: uint8_t, 标记列表末尾，固定为0xFF

            type struct zlentry { //并非实际存储的数据结构，只是便于操作而设
                unsigned int prevrawlensize;
                unsigned int prerawlen;
                unsigned int lensize;
                unsigned int len;
                unsigned int headersize;
                unsigned char encoding;
                unsigned char *p;
            } zlentry;

            例：
            *  [0f 00 00 00] [0c 00 00 00] [02 00] [00 f3] [02 f6] [ff]
            *        |             |          |       |       |     |
            *     zlbytes        zltail    entries   "2"     "5"   end


    双向链表(quicklist):
        结构(quicklist.h)：
            typedef struct quicklistNode {
                struct quicklistNode *prev;
                struct quicklistNode *next;
                unsigned char *zl;      //指向一个ziplist
                unsigned int sz;             /* ziplist size in bytes */
                unsigned int count : 16;     /* count of items in ziplist */
                unsigned int encoding : 2;   /* RAW==1 or LZF==2 */
                unsigned int container : 2;  /* NONE==1 or ZIPLIST==2 */
                unsigned int recompress : 1; /* was this node previous compressed? */
                unsigned int attempted_compress : 1; /* node can't compress; too small */
                unsigned int extra : 10; /* more bits to steal for future usage */
            } quicklistNode;

            typedef struct quicklist {
                quicklistNode *head;
                quicklistNode *tail;
                unsigned long count;        /* total count of all entries in all ziplists */
                unsigned long len;          /* number of quicklistNodes */
                int fill : 16;              /* fill factor for individual nodes */
                unsigned int compress : 16; /* depth of end nodes not to compress;0=off */
            } quicklist;


    类型与编码：
        类型                编码
        REDIS_STRING        REDIS_ENCODING_INT
        REDIS_STRING        REDIS_ENCODING_EMBSTR   //EMB编码的SDS
        REDIS_STRING        REDIS_ENCODING_RAW  //SDS
        REDIS_LIST          REDIS_ENCODING_ZIPLIST
        REDIS_LIST          REDIS_ENCODING_LINKEDLIST
        REDIS_HASH          REDIS_ENCODING_ZIPLIST
        REDIS_HASH          REDIS_ENCODING_HT
        REDIS_SET           REDIS_ENCODING_INTSET
        REDIS_SET           REDIS_ENCODING_HT
        REDIS_ZSET          REDIS_ENCODING_ZIPLIST
        REDIS_ZSET          REDIS_ENCODING_SKIPLIST

        所有redis对象均为robj类型：
        typedef struct redisObject {
        unsigned type:4;
        unsigned encoding:4;
        unsigned lru:LRU_BITS; /* LRU time (relative to global lru_clock) or
                                * LFU data (least significant 8 bits frequency
                                * and most significant 16 bits access time). */
                                //空转时间，通过object idletime访问
        int refcount;       //引用计数， 通过object refcount查询
        void *ptr;  //指向承载具体数据的地址
        } robj;

        //类型定义(server.h:466)
        #define OBJ_STRING 0    /* String object. */
        #define OBJ_LIST 1      /* List object. */
        #define OBJ_SET 2       /* Set object. */
        #define OBJ_ZSET 3      /* Sorted set object. */
        #define OBJ_HASH 4      /* Hash object. */

        //编码定义server.h:586)
        #define OBJ_ENCODING_RAW 0     /* Raw representation */
        #define OBJ_ENCODING_INT 1     /* Encoded as integer */
        #define OBJ_ENCODING_HT 2      /* Encoded as hash table */
        #define OBJ_ENCODING_ZIPMAP 3  /* Encoded as zipmap */
        #define OBJ_ENCODING_LINKEDLIST 4 /* No longer used: old list encoding. */
        #define OBJ_ENCODING_ZIPLIST 5 /* Encoded as ziplist */
        #define OBJ_ENCODING_INTSET 6  /* Encoded as intset */
        #define OBJ_ENCODING_SKIPLIST 7  /* Encoded as skiplist */
        #define OBJ_ENCODING_EMBSTR 8  /* Embedded sds string encoding */
        #define OBJ_ENCODING_QUICKLIST 9 /* Encoded as linked list of ziplists */
        #define OBJ_ENCODING_STREAM 10 /* Encoded as a radix tree of listpacks */



    列表：
        /*
        当列表中所有字符串的长度小于list-max-ziplist-value(默认64)且
        列表元素数量小于list-max-ziplist-entries(默认512)时，
        */旧设计
        新版本一律采用quicklist(linkedlist).

    哈希：
        当哈希表中所有键值的字符串长度均小于hash-max-ziplist-value(默认64)且
        键值对数量小于hash-max-ziplist-entried(默认512)时，
        采用ziplist存储，否则用hashtable。

    集合：
        当集合中所有元素都是整数且元素数量不超过set-max-intset-entries(默认512)时，
        采用intset，否则用hashtable。

    有序集合：
        同时使用skiplist/ziplist和hashtable。

        当元素数量小于zset-max-ziplist-entried(默认128)且所有元素长度小于zset-max-ziplist-value(默认64)，
        则使用ziplist，否则用skiplist。

        typedef struct zset {
            dict *dict;
            zskiplist *zsl;
        } zset;

        字典 -> 保证查找复杂度O(n)
        跳跃表 -> 保证范围性操作(ZRANGE/ZRANK等)复杂度为O(n)

        
    数据库总体设计结构：
        struct RedisServer {
            ...
            redisDb *db;    //初始化：server.db = zmalloc(sizeof(redisDb)*server.dbnum);
            ...
        };  ->  struct RedisServer server (全局变量server)

        typedef struct redisDb {
            ...
            dict *dict; // (char *)server.db.dict.ht[0].table[0].key -> key名 (由于是定长字符串，所以不用sds存储)
            dict *expires;  //过期键， 键名: 过期毫秒时间戳
            ...
        }

    过期：
        EXPIRE -> PEXPIRE -> PEXPIREAT
        EXPIREAT -> PEXPIREAT

        删除策略：
            定时删除  ->  设定定时器timer，到期时删除键值 (CPU不友好)
            惰性删除  ->  获取键值时检查过期时间 (内存不友好)   -> db.c:expireIfNeeded(redisDb *db, robj *key)
            定时删除  -> 每个一段时间检查数据库中的键是否过期 (上述2种的一种折衷，合理指定删除频率和执行时长) -> expire.c:activeExpireCycle()

        生成RDB时，过期键会被忽略。
        载入RDB时，主服务器下会忽略，从服务器下会全部载入。
        而当过期键被惰性删除或定时删除后，redis会向AOF文件追加一条DEL命令。
        因此AOF重写时过期键也不会写入到redis中。

        复制，主服务器碰到过期键会删掉自己的并向所有从服务器发送一条DEL命令，
        而从服务器执行读命令时即时碰到过期键也不会删除该过期键，而是正常返回该键值。(保证主从一致)


    事件模型：
        




    复制：
        分为同步和命令传播。
        同步：从服务器执行SLAVEOF时会先向主服务器发SYNC命令，主服务器收到后执行BGSAVE，并开设一个缓冲区记录
            当前开始的所有写命令，主服务器将BGSAVE生成的RDB文件以及缓冲区里的写命令发送给从服务器，以此来实现
            主从同步。
        命令传播： 主从同步后主服务器收到写命令会同步给从服务器实现主从一致。

        redis2.8以前的主从同步会传送包含主服务器上所有数据的RDB文件，占用大量CPU/网络/磁盘IO。
        redis2.8开始用PSYNC代替SYNC命令，并增加了部分重同步(相对于完整重同步而言)

        部分重同步由3个功能构成：
            - 主从服务器的复制偏移量 (replication offset)
            - 主服务器的复制积压缓冲区 (replication backlog) -> 默认1MB的FIFO缓冲区，所有写命令在这里都有一份副本
            - 服务器的运行ID (run ID)
        ⭐复制积压缓冲区的大小通常由second * write_size_per_second决定。
           second为从服务器断线重连所需时间，write_size_per_second为主服务器每秒产生的写命令长度(协议长)。

        完整重同步： 从服务器未复制过或者执行过SLAVEOF NO ONE时，则会发送PSYNC ? -1给主服务器
        部分重同步： 如果从服务器复制过，则发送PSYNC <runid> <offset>给主服务器，runid为上次复制主服务器的运行id，
                    offset是从服务器当前复制偏移量。
                    主服务器有三种回复：
                        - 回复+ FULLRESYNC <runid> <offset>，执行完整重同步，参数分别为主服务器的运行id和复制偏移量
                        - 回复+CONTINUE表示将执行部分重同步
                        - 回复-ERR表面主服务器版本低于2.8无法识别PSYNC命令，从服务器需要重发SYNC命令执行完整重同步

        复制的流程：
            建立TCP连接 -> 发送PING -> 发送AUTH xx(如果从服务器设置masterauth选项的话) -> 发送端口信息(REPLCONF listening-port <port>)
                 -> 同步(PSYNC) -> 命令传播
            ⭐认证时从服务器设置了masterauth且主服务器设置了requirepass或者两者都没有设置才会认证成功
            ⭐命令传播阶段从服务器会每秒发送一次REPLCONF ACK <replication_offset>给主服务器通告复制偏移量


    Sentinel：
        Sentinel是一个运行在特殊模式下的redis服务器，只能执行以下几个命令：
            - ping
            - sentinel
            - subscribe
            - unsubscribe
            - psubscribe
            - punsubscribe
            - info
        
        Sentinel通过向主服务器发送INFO命令来获得从属服务器的地址信息，并为这些从服务器创建相应的实例，同时创建连向
        这些从服务器的命令连接和订阅连接。

        Sentinel每10秒一次向监视的主从服务器发送INFO命令。当主服务器处于下线状态或正在对主服务器进行故障转移时，向
        从服务器发送INFO的频率变为1秒1次。

        对于监视同一个主从服务器的多个Sentinel来说，它们会每2秒1次向被监视服务器的__sentinel__:hello频道发送消息来
        向其它Sentinel宣示自己的存在。

        Sentinel也会从__sentinel__:hello频道里接收其它Sentinel的信息并依此创建相应的实例和命令连接。

        Sentinel与主从服务器之间会创建命令连接和订阅连接，而Sentinel之间只会创建命令连接。

        Sentinel以每秒1次的频率向所有实例(主从服务器/其它Sentinel)发送PING命令，根据回复判断是否在线。连续收到无效
        回复时会判断它为主观下线。

        当一个主服务器被判断为主观下线时，Sentinel会向其从服务器发送询问是否同意该主服务器已进入主观下线状态。

        收集到足够投票时则判断主服务器为客观下线状态并发起针对该主服务器的故障转移操作。


    集群：
        节点通过握手将其它节点加入到自己当前集群中。 (CLUSTER MEET命令)

        集群中共有16384个槽，所有槽都可以被处理时集群才属于上线状态。通过CLUSTER ADDSLOTS可以将某个槽指定给某节点处理。

        节点收到请求时会先查看是否属于自己处理，否则发出MOVED错误指引请求正确的节点。

        redis-trib负责集群的重新分片工作，重新分片的关键是将某个槽的所有键值对迁移到另一个槽。
    
    发布与订阅：
        struct redisServer {
            ...
            dict *pubsub_channels;  /* Map channels to list of subscribed clients */
            list *pubsub_patterns;  /* A list of pubsub_patterns */
        };
        redisServer里的pubsub_channels里存放了所有订阅的对应关系，对应关系为channel <-> 订阅该channel的client列表。
        如"news-it" : client-9 -> client-23 -> client-7

        typedef struct pubsubPattern {
            client *client;
            robj *pattern;
        } pubsubPattern;

        而pubsub_patterns里是一个pubsubPattern类型的链表，链表节点是client和订阅模式，如client-12，"news-*"。

        struct client {
            ...
            dict *pubsub_channels;
            list *pubsub_patterns;
        }；

        此外，每个client结构里也维护了一份订阅频道和订阅模式，和server里的信息保持一致，sub/unsub的时候先从client里
        操作然后决定是否继续下去。
        如：    if (dictAdd(c->pubsub_channels,channel,NULL) == DICT_OK) { ... //sub时
                if (dictDelete(c->pubsub_channels,channel) == DICT_OK) { ... //ubsub时


    事务:
        /* Client MULTI/EXEC state */
        typedef struct multiCmd {
            robj **argv;
            int argc;
            struct redisCommand *cmd;
        } multiCmd;

        typedef struct multiState {
            multiCmd *commands;     /* Array of MULTI commands */
            int count;              /* Total number of MULTI commands */
            int cmd_flags;          /* The accumulated command flags OR-ed together.
                                    So if at least a command has a given flag, it
                                    will be set in this field. */
            int minreplicas;        /* MINREPLICAS for synchronous replication */
            time_t minreplicas_timeout; /* MINREPLICAS timeout as unixtime. */
        } multiState;

        typedef struct client {
            int flags;              /* Client flags: CLIENT_* macros. */
            multiState mstate;      /* MULTI/EXEC state */
            list *watched_keys;     /* Keys WATCHED for MULTI/EXEC CAS */
        }

        redis用FIFO队列实现事务队列,事务的命令队列为client.mstate.commands[0]/commands[1]...

        typedef struct redisDb {
            ...
            dict *watched_keys;         /* WATCHED keys for MULTI/EXEC CAS */
        }

        redisDb的watched_keys存储了监视keys与client的对应关系,当对dict的keys有改动的时候,会遍历对应的client链表
        并通过touchWatchedKey()函数将client结构里的flag变量加入CLIENT_DIRTY_CAS,当EXEC时检查对应client的flag里
        是否有该CLIENT_DIRTY_CAS,有的话则放弃执行并返回nil.

        关于redis事务的ACID:
            - 原子性, 事务中的命令要么全执行,要么一个都不执行.redis满足该特性但不支持事务回滚.
              如: 命令队列入队时发现命令参数数量不对(GET后面没接key)则全部放弃执行.
                  执行过程中发现命令和参数类型不符合(对字符串RPUSH这种),则会跳过继续执行.
            - 一致性, 执行前数据库是一致性的,则执行后依旧满足一致性,
              redis下有谨慎的错误检测,如下面三种情况:
                - 入队错误, 比如输入了不存在的命令. 此时所有命令都不会被执行
                - 执行错误, 这个无法在入队检测时被发现,执行时会跳过它
                - 服务器宕机, 无论在非持久/RDB/AOF哪种模式下都不会影响一致性
            - 隔离性, 多个事务并发执行时不会互相影响. redis是单进程单线程操作,各事务处于串行执行,因此满足隔离性.
            - 耐久性, 事务执行完毕时事务的执行结果会被保存到永久介质里。
              redis下只有appendfsync为always时，每次执行完指令会调用sync命令同步到磁盘里，此时才算满足耐久性。

    二进位制数组：
        SETBIT/GETBIT/BITCOUNT/BITOP

        位数组用sds来表示，robj.ptr -> sds{len:..., alloc:..., flags:..., buf[0],buf[1]...buf[N]} //其中最后的buf[N]为'\0'.

        汉明重量：位数组中非0位的数量。
        redis采用查表和variable-precisionSWAR结合的方法来计算BITCOUNT。
            - 先每28字节用SWAR算法计算，最后剩余的不足28字节用查表来得出。
    

    慢查询日志：
        slowlog-log-slower-than X表示执行时间大于X微秒的命令将会被记录。 //默认10000即10毫秒
        slowlog-max-len指定服务器保存的慢查询日志的上限条数。 //默认128

        slowlog get获取所有慢查询日志，slowlog len返回慢查询日志数量。

