refer:  https://www.jianshu.com/p/75edcc05acfd

# 1. 概念

单点登录，即Single Sign On，在多个应用系统中，只需要登录一个就可以访问其它互相信任的系统。

# 2. 同域下的单点实现

例如：有一个域名app.com，下属两个应用x.app.com和y.app.com，单点登录地址为sso.app.com.

    - step1：登录x.app.com时发现用户没有凭证，让其转到sso.app.com登录
    - step2：客户端访问sso.app.com登录成功后获得凭证，写到sso.app.com的cookie下
    - step3. 由于cookie不能跨域所以客户端再次访问x.app.com时无法携带之前获得的凭证
             并且即使x.app.com得到了这个凭证也需要和sso.app.com共享session才能认证
    - step4. 设定cookie的域时只能设置当前域和顶域，故可设置凭证cookie的域为app.com，
            这样访问x.app.com可以带上这个cookie，另外x.app.com/y.app.com/sso.app.com
            还需共享session。

# 3. 不同域的单点登录
    术语：  TGC=Ticket Granting Ticket
            ST=Service Ticket

    when accessing app.example.com:
    - client->app: GET https://app.example.com
    - client<-app: 302 Location: https://cas.example.com/cas/login?service=https%3A%2F%2Fapp.example.com
    - client->cas: GET https://cas.example.com/cas/login?service=https%3A%2F%2Fapp.example.com
    - client<-cas: Login Form
    - client->cas: POST https://cas.example.com/cas/login?service=https%3A%2F%2Fapp.example.com with username+password
    - client<-cas: authenticate success, create CASTCG cookie which contains TGT(Ticket Granting Ticket), the session key for SSO session.
                    Set-Cookie: CASTGT=TGT-2586527
                    302 Location=https://app.example.com/?ticket=ST-1234567
    - client->app: GET https://app.example.com/?ticket=ST-1234567
    - app->cas   : GET https://cas.example.com/serviceValidate?service=https%3A%2F%2Fapp.example.com&ticket=ST-1234567
    - cas->app   : 200 validate OK
    - client<-app: 302 Location: https://app.example.com
                   Set-Cookie: JSESSIONID=ABC111111
    - client->app: GET https://app.example.com
                    Cookie: JSESSIONID=ABC111111
    - client<-app: 200 Content      //validate session cookie JSESSIONID success


    then accessing app2.example.com:
    - client->app2: GET https://app2.example.com
    - client<-app2: 302 Location: https://cas.example.com/cas/login?service=https%3A%2F%2Fapp2.example.com
    - client->cas:  GET https://cas.example.com/cas/login?service=https%3A%2F%2Fapp2.example.com
                    Cookie: CASTGT=TGT-2586527
    - client<-cas:  validate CASTGT ok then no login is required
                    302 Location: https://app2.example.com/?ticket=ST-9876543
    - client->app2: GET https://app2.example.com/?ticket=ST-9876543
    - app2->cas:    GET https://cas.example.com/serviceValidate?service=https%3A%2F%2Fapp2.example.com&ticket=ST-9876543
    - app2<-cas:    200 validate OK
    - client<-app2: 302 Location: https://app2.example.com/
                    Set-Cookie: MOD_AUTH_CAS_S=XYZ330012
    - client->app2: GET https://app2.example.com
                    Cookie: MOD_AUTH_CAS_S=XYZ330012
    - client<-app2: 200 Content      //validate cookie MOD_AUTH_CAS_S




