
etcd主流应用场景：
    1. 服务发现(service discovery)
        - 强一致/高可用的kv键值存储
        - 注册服务和监控服务健康状态
        - 正常查找和连接服务的机制
    
    2. 消息发布与订阅
        应用启动时主动get一次配置，然后注册一个watcher并等待，保证配置变动时能得到通知。
    

    3. 分布式协调
        - 分布式锁
            - 保持独占， etcd通过CAS(CompareAndSwap)可以实现分布式场景下多个客户端创建同一个目录时只有一个会创建成功即视为获得了锁
                如：c1, c2, c3同时争夺一个lockA， 最后c2成功则etcd里表现为/lockA=c2
            - 控制时序， 多客户端取一个锁时获取锁的顺序是全局唯一并且决定了执行顺序。
                分布式队列，如c1, c2, c3同时获取orderA，则可能表现为/orderA/1=c2, /orderA/2=c3, /orderA/3=c1,即执行顺序为c2,c3,c1

    4. 集群监控
        利用watcher机制监视节点，并且节点可以设置TTL，每隔一定时间更新key表明自己还在存活。
    
    5. leader竞选
        利用CAS机制选出leader避免一些重复性的劳动(CPU/IO操作)。
        比如搜索系统中的全量索引，先用etcd选出leader节点算出索引再分发到其它节点。


etcd架构：

request ->  HTTP Server <->  Store <-> Raft  
                        <->  Raft   <-> WAL
                                    <-> SnapShot

    - HTTP Server： 用于处理用户发送的 API 请求以及其它 etcd 节点的同步与心跳信息请求。
    - Store：用于处理 etcd 支持的各类功能的事务，包括数据索引、节点状态变更、监控与反馈、事件处理与执行等等，
        是 etcd 对用户提供的大多数 API 功能的具体实现。
    - Raft：Raft 强一致性算法的具体实现，是 etcd 的核心。
    - WAL：Write Ahead Log（预写式日志），是 etcd 的数据存储方式。除了在内存中存有所有数据的状态以及节点的索引以外，etcd 就通过 WAL 进行持久化存储。    WAL 中，所有的数据提交前都会事先记录日志。Snapshot 是为了防止数据过多而进行的状态快照；Entry 表示存储的具体日志内容。

部分术语： 
    snapshot：etcd 防止 WAL 文件过多而设置的快照，存储 etcd 数据状态。
    Proxy：etcd 的一种模式，为 etcd 集群提供反向代理服务。
    Leader：Raft 算法中通过竞选而产生的处理所有数据提交的节点。
    Follower：竞选失败的节点作为 Raft 中的从属节点，为算法提供强一致性保证。
    Candidate：当 Follower 超过一定时间接收不到 Leader 的心跳时转变为 Candidate 开始竞选。
    Term：某个节点成为 Leader 到下一次竞选时间，称为一个 Term。
    Index：数据项编号。Raft 中通过 Term 和 Index 来定位数据。    
