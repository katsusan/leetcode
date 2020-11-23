# 0.Requirements

语言要求：GoLang，Java（ 主修）
网络编程：Netty，Mina
    多线程，并发，Concurrent包
    JVM调优，GC，熟悉诊断工具
    设计模式，JDK，异步同步通讯
熟悉常用中间件：
    关系型数据库（MySql）；
    NoSql(Redis，Hbase，MongoDB, levelDB)；
    分布式服务框架（Spring Cloud，gRPC， Dubbo）；
    分布式消息队列（Kafka，RabbitMQ，ActiveMQ）；
    搜索引擎（ElasticSearch）；
    分布式工具（ZooKeeper）；
    分布式配置中心（Diamond）；
环境和工具：
    Linux，GIT，IDEA
    docker， K8S
    DevOps
基础：
网络：TCP，HTTP，HTTP2
操作系统：
数据结构及算法：
    分布式理论（存储，跟踪，一致性算法）
    缓存技术

# 1.跨域相关：

浏览器的同源策略：
    满足"域名+IP+端口"相同则为同源请求，
    同源限制以下几项：
        - Cookie、LocalStorage、IndexedDB 等存储性内容
        - DOM 节点
        - AJAX 请求发送后，结果被浏览器拦截了
    但有三个标签允许跨域加载：
        - <img src=XXX>
        - <link href=XXX>
        - <script src=XXX>
    跨域是浏览器行为，服务端会收到请求并返回结果。


解决方案：

    1. 利用script标签的未被限制策略，可以由此请求服务端数据，但只支持get，且易受XSS攻击。
    2. CORS。需浏览器和服务端同时支持，服务端设置Access-Control-Allow-Origin为允许请求的域。


# 2.Web安全相关：

1. XSS
    原理是在网页里嵌入可执行的脚本代码实现某种目的。
    A. 非持久型XSS(反射型)   ->  在URL中放入可执行脚本的参数
        应对: 网页中渲染内容保证来自服务端.
                前端对渲染内容做必要的转义.并少使用eval等直接执行字符串的方法.
    B. 持久型XSS(存储型)     ->  一般存在于form表单等功能,比如评论/文章等里面写入特定js代码.

    通用防御方法:
    CSP  ->  设定HTTP Header中的Content-Security-Policy / 设定meta标签的方式.
                参考：https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/
    字符转义 ->  str.replace(/&/g, '&amp;'), str.replace(/</g, '&lt;'), ...
    Cookie设置httponly   ->  服务器端响应首部 Set-Cookie: key=value; HttpOnly, js就无法用document.cookie取得cookie

2. CSRF
    Cross Site Request Forgery, 跨站请求伪造。
    应对方法：
    SameSite    ->  对Cookie设置SameSite属性表示不会随跨域发送，但浏览器支持度有限
    Referer Check   ->  通过检测refer头判断是否由当前页面来防范，但https页面里向http链接请求时通常不会发送refer头
    CSRF token  ->  服务端为每个session维护一个token(通常是变动的),对每次请求都验证请求中携带的token与服务端的是否一致。


# 3.分布式

分布式：
    事务的ACID特性:  Atomicity ->  原子性
                    Consistency ->  一致性
                    Isolation   ->  隔离性
                    Durability  ->  持久性
    
    分布式CAP：      C   ->  Consistency
                    A   ->  Available
                    P   ->  Partition

    BASE:   (对CAP让步) 基本可用 + 软状态 +最终一致
            Basically Available + Soft-state + Eventually Consistency


# 4.微服务

SOA ：Service-oriented-Architecture 面向服务的架构

服务治理的几大问题：
    -访问权限   (比如限定用户A访问接口X)
    -版本控制
    -服务时效/次数限制
    -性能措施 (高TPS时的引流or缓存等)
    -跨平台