
1.知乎项目结构

├── bin              	--> 构建生成的可执行文件
├── cmd              	--> 各种服务的 main 函数入口（ RPC、Web 等）
│   ├── service 
│   │    └── main.go
│   ├── web
│   └── worker
├── gen-go           	--> 根据 RPC thrift 接口自动生成
├── pkg              	--> 真正的实现部分（下面详细介绍）
│   ├── controller
│   ├── dao
│   ├── rpc
│   ├── service
│   └── web
│   	├── controller
│   	├── handler
│   	├── model
│   	└── router
├── thrift_files     	--> thrift 接口定义
│   └── interface.thrift
├── vendor           	--> 依赖的第三方库（ dep ensure 自动拉取）
├── Gopkg.lock       	--> 第三方依赖版本控制
├── Gopkg.toml
├── joker.yml        	--> 应用构建配置
├── Makefile         	--> 本项目下常用的构建命令
└── README.md


- bin：构建生成的可执行文件，一般线上启动就是 `bin/xxxx-service`
- cmd：各种服务（RPC、Web、离线任务等）的 main 函数入口，一般从这里开始执行
- gen-go：thrift 编译自动生成的代码，一般会配置 Makefile，直接 `make thrift` 即可生成（这种方式有一个弊端：很难升级 thrift 版本）
- pkg：真正的业务实现（下面详细介绍）
- thrift_files：定义 RPC 接口协议
- vendor：依赖的第三方库

pkg下：

pkg/
├── controller    	
│   ├── ctl.go       	--> 接口
│   ├── impl         	--> 接口的业务实现
│   │	└── ctl.go
│   └── mock         	--> 接口的 mock 实现
│     	└── mock_ctl.go
├── dao           	
│   ├── impl
│   └── mock
├── rpc           	
│   ├── impl
│   └── mock
├── service       	--> 本项目 RPC 服务接口入口
│   ├── impl
│   └── mock
└── web           	--> Web 层（提供 HTTP 服务）
    ├── controller    	--> Web 层 controller 逻辑
    │   ├── impl
    │   └── mock
    ├── handler       	--> 各种 HTTP 接口实现
    ├── model         	-->
    ├── formatter     	--> 把 model 转换成输出给外部的格式
    └── router        	--> 路由





*****************************************MyProject***********************************************
main()
    |
ginadmin.Init()
    |
InitLogger()
InitCaptcha()
InitObject()    ->  Object{*casbin.Enforcer, auth.Auther, model.Common, bll.Common} ->  InitJWTAuth(), InitStore(), 
InitData()
InitWeb()
InitHttpServer()


# centaur

current mod:
    github.com/Katsusan/centaur

    - config
        "github.com/Katsusan/centaur/internal/config"
        
web:
    github.com/gin-gonic/gin

config:
    github.com/spf13/viper

orm:
    github.com/jinzhu/gorm

log:
    log "github.com/sirupsen/logrus"


****************************************************API*******************************************************

GET /api/v1/login/captchaid: 
    =>  "captcha_id":"w9HT2bHXdAmtUEtq5f7v"

GET /api/v1/login/captcha?captcha_id=w9HT2bHXdAmtUEtq5f7v:
    =>  xxx.png (验证码图片)

POST /api/v1/login:
    {"user_name":"root",
    "captcha_code":"3254",
    "captcha_id":"w9HT2bHXdAmtUEtq5f7v",
    "password":"6351623c8cef86fefabfa7da046fc619"}
    =>  access_token: "~" (jwt token)
        expires_at: 1565958511
        token_type: "Bearer"






