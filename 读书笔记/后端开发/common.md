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
        - `<img src=XXX>`
        - `<link href=XXX>`
        - `<script src=XXX>`
    跨域是浏览器行为，服务端会收到请求并返回结果。


解决方案：

    1. 利用script标签的未被限制策略，可以由此请求服务端数据，但只支持get，且易受XSS攻击。
    2. CORS。需浏览器和服务端同时支持，服务端设置Access-Control-Allow-Origin为允许请求的域。


# 2.Web安全相关：

## 2.1 XSS
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

## 2.2 CSRF
    Cross Site Request Forgery, 跨站请求伪造。
    应对方法：
    SameSite    ->  对Cookie设置SameSite属性表示不会随跨域发送，但浏览器支持度有限
    Referer Check   ->  通过检测refer头判断是否由当前页面来防范，但https页面里向http链接请求时通常不会发送refer头
    CSRF token  ->  服务端为每个session维护一个token(通常是变动的),对每次请求都验证请求中携带的token与服务端的是否一致。

## 2.3 

## 2.3.1 认证(Authentication)
用户身份的验证，通常如：
- 用户名密码验证
- 邮箱发送验证码
- 手机验证码认证
- 第三方登录

## 2.3.2 授权(Authorization)
授予用户访问某些资源的权限：
- 比如手机应用请求的权限(麦克风，相机，位置信息等)
- 访问某个API的权限

授权的方式有：cookie，session，token，OAuth，OpenID，等等。

## 2.3.3 凭证(Credential)
实现认证和授权时标记用户身份的媒介，比如：
- 实体凭证，如个人身份证，银行卡，护照等
- 用户登录后的凭证，如颁发的token。

## 2.3.4 Cookie
HTTP属于无状态协议，为了识别当前请求的身份信息，需要一种手段来保存会话上下文，这种手段就是Cookie。

cookie是存储在浏览器上的一段文本，可以通过HTTP请求头的Set-Cookie和Set-Cookie2(已弃用)字段来设置。
会在下一次向同一个服务器发送请求时放在Cookie请求头中，以便服务器可以识别客户端。

cookie只能在同一个域名(domain属性)下使用，不能跨域使用。

Cookies are mainly used for three purposes:
- Session management
    Logins, shopping carts, game scores, or anything else the server should remember
- Personalization
    User preferences, themes, and other settings
- Tracking
    Recording and analyzing user behavior

Set-Cookie用于服务器返回到用户端时设置cookie。

Server -> Client: 
    SetCookie: usr=yax; path=/; domain=.example.com; expires=Thu, 01-Jan-2021 00:00:01 GMT; secure; httponly
Client -> Server: 
    GET / HTTP/1.1
    Host: www.example.com
    Cookie: usr=yax;

document.cookie属性可以获取当前域名下的所有cookie。
typeof document.cookie → 'string'
document.cookie = 'foo=bar;path=/;max-age='+5*60+';';


// https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies   
cookie属性：
- name：cookie名称
- content：设置cookie的值
- domain()：设置cookie所在的域名，默认当前域名。subdomain也可以访问，比如mozilla.org时也可以被developer.mozilla.org访问。
- path：设置cookie所在的路径，假如设为/abc,则/abc/def也可以使用，默认为/
- maxAge: 设置cookie的有效期，单位为秒，默认为-1，表示关闭浏览器清除cookie,浏览器看到maxAge会忽略expires。
- expires: 绝对过期时间。
- secure：设置cookie是否只能带在https下的请求(除了localhost)里，默认为false。
- httpOnly：设置是否可以被document.cookie访问，带httponly的cookie无法被上述JavaScript API访问，默认为false。这个减少了XSS攻击的可能性。
- samesite: 设置cookie的同源策略，可选Strict/Lax(default)/None。
            Strict代表只在同源请求里携带该cookie，Lax类似Strict，除了用户被navigate到origin site时也可以携带。
            None代表浏览器在cross-site和same-site请求里都会带该cookie，此时secure也必须被设置，即`SameSite:None;Secure`。


## 2.3.5 session
session是一种记录了客户端会话状态的机制。存在服务器断，在客户端认证成功后用Set-Cookie设置session-id，
每次请求时，客户端会把session-id放在cookie中，以便服务器可以识别客户端。

flow：
- 客户端第一次请求服务端时，提交认证信息，成功的话服务器为其创建对应的session。
- 请求返回时将此session-id放在Set-Cookie下。
- 浏览器收到服务端响应后将Set-Cookie中的session-id放入浏览器的cookie中。
- 浏览器再次请求服务器时，将cookie中的session-id放入请求头中，服务器检查session-id是否存在，存在则继续，不存在则重定向到登录或返回失败。

## 2.3.6 cookie和session的区别
cookie是存储在浏览器上的一段文本，可以通过HTTP请求头的Set-Cookie和Set-Cookie2(已弃用)字段来设置。
session是存放在服务端用来记录客户端会话状态的数据。 

cookie大小最多4KB，session大小由服务端决定。

## 2.3.7 Token
访问API所需要的凭证。
简单Token的组成：uid(用户唯一标识)，time(当前时间戳)，signature(token的有效数据的哈希值)。
特点：
- 服务端无状态化，可扩展性好
- 支持移动端设备
- 安全
- 支持跨程序调用

flow：
- 客户端提交认证信息，服务端返回token。
- 服务端收到请求验证用户名密码。
- 验证成功后服务端签发一个token并把token放在Set-Cookie中。
- 客户端收到响应后将token放入cookie/localStorage中。
- 客户端每次向服务端请求都要携带token，服务端验证token是否有效。
- 服务端收到请求时验证token是否有效，有效则返回请求所需数据。

token用计算和解析token的时间换取了session的空间，减轻了用session时查询数据库的压力。

## 2.3.8 refresh token
refresh token是用来刷新access token的token，它的组成与access token相同，但是过期时间比access token长。
没有refresh token的话也可以刷新access token，但每次刷新都需要用户重新认证，服务端也要做相应的处理。
有了refresh token的话，可以不用重新认证，直接刷新access token。

flow:
- 客户端提交认证信息，服务端返回2个token, access token有效期：1周；refresh token有效期：1个月。
- 客户端请求时携带access token，服务端验证access token是否有效，有效则返回请求所需数据，无效则返回invalid token。
- 客户端携带refresh token，服务端验证refresh token是否有效，有效则返回新的access token，无效则需重新登录。

refresh token和其过期时间存储在数据库中，只有申请新的acccess token时才会验证refresh token，因此不会有太大的查询压力。
也不需要像session一样需要保存在内存中应对大量的请求。

## 2.3.9 token vs session
session是记录服务端和客户端会话状态的机制，使服务端有状态化，可以记录会话信息。
token是访问资源的凭证，服务端无状态化，扩展性好。

在验证方面token的安全性好，token可以防止重放攻击(replay attack),而session只能依赖传输层加密了。

## 2.3.9 JWT
JSON Web Token，简称JWT，是目前最流行的跨域认证方案。标准定义在rfc7519(https://tools.ietf.org/html/rfc7519)。

流程：认证成功 → 服务端用只有自己知道的secret对header&payload签名生成JWT → 客户端通常将返回的JWT存在local storage中，也可以存于cookie。

形式： 
- HTTP header：`Authorization: Bearer <token>`
- 跨域的时候放在post的body中
- 放在url参数中`/user?token=<token>`

JWT不用cookie的时候，可以使用任意域名的api不用担心跨域问题，

example:
    https://github.com/yjdjiayou/jwt-demo
    https://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html

## 2.3.10 JWT vs token
都是资源访问的令牌/都可以记录用户信息/都可以使服务器无状态化

区别：
- Token：服务端收到token时需要查询数据库验证用户信息
- JWT：服务端收到JWT时只需要解密就可以验证用户信息，不需要查询数据库

## 2.3.11 常见的前后端鉴权
- Session-Cookie
- Token(JWT,SSO等)
- OAuth2(开放授权)

## 2.3.12 常见问题
- Cookie：
	+ 存在客户端，容易被客户端篡改，使用前需要验证
	+ 不要存敏感数据，如密码，余额等
	+ 使用httponly提高安全性
	+ 设置正确的domain和path减少传输
	+ cookie无法跨域
	+ 通常来讲单个网站最多20个cookie，每个cookie大小<4kb，浏览器允许<300个cookie
	+ 移动端对cookie支持不好，而session基于cookie实现，因此通常用token

- Session：
	+ 存储在服务端，需要维护一个session数据库，当用户量变大时会导致性能问题
	+ session共享问题，集群部署时维护一致性
	+ session-id存在cookie中，如果用户禁用了cookie或不支持cookie那么需要把session-id放在url参数中

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