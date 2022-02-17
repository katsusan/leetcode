from https://yeasy.gitbooks.io/docker_practice/content/

example：
    docker run --detach=true --name=redisx -p 6379:6379 redis  (--detach 后台)
    docker run -it -p 6379:6379 redis /bin/bash (-it 前台,  /bin/bash 进入docker内终端)
    docker exec -it redisx /bin/bash


# 1. 命名空间namespace
    linux2.6之后的内核加入了一种资源隔离机制，可以将内核中支持的系统资源比如进程/用户账户/文件系统/网络
    分属为特定的命名空间，每个命名空间下的资源对于其它的命名空间资源是透明不可见的，因此在不同命名空间下
    可能出现相同pid的进程。

    以linux4.15.14内核为例，有8种不同的命名空间属性。用ll /proc/[pid]/ns可以看到某个pid下的命名空间属性。
    ➜  ~ ll /proc/1/ns
    total 0
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 cgroup -> cgroup:[4026531835]
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 ipc -> ipc:[4026531839]
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 mnt -> mnt:[4026531840]
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 net -> net:[4026531993]
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 pid -> pid:[4026531836]
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 pid_for_children -> pid:[4026531836]
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 user -> user:[4026531837]
    lrwxrwxrwx 1 root root 0 Sep 11 18:21 uts -> uts:[4026531838]

    后面的编号代表命名空间的编号，如果不启动容器的化一般各个进程的这个值一样。

    - uts：主机名隔离。指每个独立容器空间内的程序可以拥有各自不同的主机名信息。
    - pid：基于进程的隔离。使得容器中的首个进程成为所在命名空间中PID为1的进程。
        pid为1的进程是所有进程的根父进程，有很多特权如屏蔽信号/接管孤儿进程。当它退出时同属一个命名空间的所有进程都会被杀死。
    - ipc：基于System V的进程通信隔离。IPC(Inter-Process Communication)是linux中标准的进程间通信方式，包括共享内存/信号量/消息队列等。
        ipc隔离使得同一个命名空间下的进程才能相互通信。
    - mnt：基于磁盘挂载点和文件系统的隔离。mnt namspace会为隔离空间创建独立的mount节点树。容器中的进程无法访问容器外的任何文件。
    - user：基于系统用户的隔离。指同一个系统用户在不同的命名空间中可以拥有不同的uid/gid。即特定命名空间中uid为0的用户并不一定是系统的root用户。
    - net：基于网络栈的隔离。允许使用者将特定网卡与特定容器中的进程运行上下文关联起来，使得同一个网卡在不同容器中显示不同的名称。


# 2. 控制组cgroup(Control Group)
    linux内核提供的一种可以限制/记录/隔离进程组所使用的物理资源(CPU/内存/磁盘IO/网络等)。
    cgroup有三个核心概念。
        - Hierarchy 层级
        - Subsystem 子系统
        - Control Group 控制组
    ~每个层级表示系统中的一组cgroup配置树，可以包含多个子系统
    ~每个子系统对应一种控制资源
    ~每个子系统下可以创建不同的控制组，控制组可以被赋予一组进程，通过指定控制组的参数来对这组进程的资源使用进行控制or记录
    即：
        层级
        ↓
        子系统(CPU子系统，内存子系统，IO子系统等)
        ↓
        控制组(控制组1，控制组2，...)
        ↓
        进程(进程1，进程2，...)

    建立cgroup的完整步骤：
        - 在文件系统上建立cgroup层级的目录
            linux内核默认为将加载的所有cgroup子系统挂载刀"/sys/fs/cgroup"下，这个目录就表示了默认的一个cgroup层级。
            通过 mount -t cgroup可以将其打印出来。实际操作cgroup时可以省去挂载层级和子系统的操作，直接在下面创建控制组。
        - 挂载文件系统并关联子系统
            可以使用系统默认的/sys/fs/cgroup而省略这步操作
            或者自己创建的话如下：
                mkdir cgroup
                mount -t tmpfs cgroup_root ./cgroup //将特殊的cgroup根设备挂载到这个目录
                mkdir cgroup/cpu
                mount -t cgroup -o cpu cgroup ./cgroup/cpu  //cpu子系统
        - 建立控制组
                mkdir cgroup/cpu/thirty_percent //限制CPU使用率30%的控制组
        - 设置控制参数
                cat cgroup/cpu/thirty_percent/cpu.cfs_quota_us  //默认为-1，表示不限制
                echo 30000 > cgroup/cpu/thirty_percent/cpu.cfs_quota_us //限制为30%
        - 将进程加入到控制组中  
                echo 1122 > cgroup/cpu/thirty_percent/tasks //将进程id为1122的进程加入到这个控制组，此时cpu的30%限制即刻生效

    进程的cgroup属性会继承，因此父进程的限制对子进程依然有效，类似于树状


# 3. Dockerfile指令
    每个dockerfile指令都会创造一层镜像。
    - FROM
        表明以哪个镜像为基础进行构建，特别地，从空白镜像构建用FROM scratch(Go应用常用)
    - RUN  
        由于每个RUN指令都会生成一层中间镜像加上Union FS对镜像层数有限制，因此可以尽量压缩在一个命令行使用。
        RUN apt-get update
        RUN apt-get install gcc
        ↓
        RUN apt-get update && apt-get install gcc
    
    - COPY --chown=<user>:<group> src dst   {chown:改变所属用户和所属组}
    - ADD --chown=<user>:<group> src dst 
        和COPY类似，但多出了自动解压缩的功能，尽量使用语义明确的COPY
    - CMD
        一般用于指定容器内主进程的启动，如CMD ["nginx", "-g", "daemon off;"]
    - ENTRYPOINT 
        与CMD类似，都是为了指定启动程序和参数，在运行时也可以替代(docker run --entrypoint)。
        但用ENTRYPOINT后就不是直接运行，而是作为参数传给ENTRYPOINT。即<ENTRYPOINT> "<CMD>"
    - ENV key value/key1=value1 key2=value2
        支持环境变量展开的命令：ADD、COPY、ENV、EXPOSE、LABEL、USER、WORKDIR、VOLUME、STOPSIGNAL、ONBUILD。
    - ARG param=value
        和ENV类似，但容器以后运行时这些环境变量是不存在的。可以用docker build --build-arg param=value来覆盖。
    - VOLUME 路径（,路径2,...）
        VOLUME /data会将/data挂载为匿名卷，任何向容器/data里写入的信息都不会记录进容器存储层。
        可以用docker run -d -v mydata:/data来覆盖，代表用mydata这个命名卷挂载到了/data的位置，替代了原本定义的匿名卷。
    - EXPOSE 端口(,端口2,...)
        只是声明应用会在哪些端口提供服务。1是便于理解端口应用提供的服务，2是运行时可以用随机端口映射docker run -P。
        这个与docker run -p <宿主端口>:<容器端口>是不同的。
    - WORKDIR 路径
        指定工作目录，之后的各层工作目录就变为WORKDIR。
    - USER 用户名(:用户组)
        改变之后各层执行CMD、RUN、ENTRYPOINT这类命令执行时的身份。
    - HEALTHCHECK 选项 CMD 命令        (>=Docker 1.12)
        告诉Docker如何判断容器内服务是否正常。和CMD、ENTRYPOINT一样只可以出现一次，多次以最后一次为准。
        选项：
            --interval=<间隔>, 默认为30s
            --timeout=<时长>, 默认30s
            --retries=<次数>，默认3次
        后面的命令返回值决定此次healthcheck成功与否，0：成功；1：失败；2：保留，不要使用这个值。
        状态有starting->healthy->unhealthy。
        以nginx为例：
            FROM nginx
            RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
            HEALTHCHECK --interval=5s --timeout=3s \
                CMD curl -fs http://localhost/ || exit 1
    - ONBUILD
        ONBUILD后面跟RUN、COPY等，它在当前镜像构建时并不会运行，只有别的镜像以它为基础构建时才会执行。

# 4. Docker操作
    + docker run
        -t  tty: Allocate a pseudo-TTY
        -i  interactive: Keep STDIN open even if not attached
        -d  detach: Run container in background and print container ID
    + docker container
        start/stop/restart
        ls (-a)
        rm
        prune   Remove all stopped containers
        logs
    + docker attach attach到容器内，exit会退出容器
    + docker exec
        -t detach
        -i interactive
    + docker export
        -o 将container导入到-o指定的文件(tar格式)
    + docker import
        file/url/- repo[:tag] 导入文件并生成指定容器


# 5. 私有仓库
    官方库：
    docker run -d -p 5000:5000 -v /opt/data/registry:/var/lib/registry --restart=always --name registry registry

    docker tag ubuntu:latest localhost:5000/ubuntu:latest
    docker push localhost:5000/ubuntu:latest    //will push this image to localhost:5000
    docker pull locahost:5000/ubuntu:latest     //will pull ubuntu:latest from localhost:5000

    以上官方库有一些缺点比如一些image删除后不会自动回收空间等等，下面是推荐的第三方库管理镜像。
    docker run -d --name nexus3 --restart=always -p 8081:8081 --mount src=nexus-data,target=/nexus-data sonatype/nexus3


# 6. 数据卷
    数据卷可以在容器之间共享和重用
    对数据卷的修改会立马生效
    对数据卷的更新不会影响镜像
    数据卷会一直存在即使容器被删除

    docker volume
        create x    创建数据卷x
        ls          显示所有数据卷
        rm
        prune
        inpect

    执行docker run的时候用--mount标记来将数据卷挂载到容器里，可以同时挂载多个。
        docker run --mount src=vol1, taget=/data1
    
    同时mount也可以将宿主机目录/文件挂载进容器中，如
        docker run --mount type=bind,src=/usr/docker/webdata,target=/webdata,readonly
    

# 7. docker网络
    -P参数代表docker会将49000-49900中随机的端口映射到容器内应用的开放端口。
    -p则是指定端口，支持的格式有：ip:hostPort:containerPort | ip::containerPort | hostPort:containerPort。
    例：
        -p 80:80            //将宿主机任意接口的80绑定到容器的80端口
        -p 127.0.0.1:80:80  //将宿主机本地网卡接口上的80绑定到容器的80端口
        -p 127.0.0.1::80    //绑定宿主机上127.0.0.1接口的任意端口到容器的80端口
    
    docker port <containerID> 用来查看容器的端口配置

    docker network create -d bridge <netname>
    docker run --network <netname>


# 8. Multi-stage多阶段构建

通常应用的编译环境与运行环境的要求并不要求一致，以Go来讲，编译镜像`FROM golang:1.x.y`一般要上百M，
而生产中运行其可执行文件并不需要golang环境本身，可以借助Dockerfile中的multi-stage来应对此问题。

以一个简单的echo服务为例，它把http请求的parameters原样返回，代码如下：

```Go
// echo.go
func main() {
        debug := os.Getenv("DEBUG")
        http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
                if err := r.ParseForm(); err != nil {
                        w.Write([]byte("incorrect params"))
                        return
                }
                for k, v := range r.Form {
                        // name = john,bob
                        w.Write([]byte(k + " = " + strings.Join(v, ",")))
                }
        })
        http.ListenAndServe(":8000", nil)
}
```

# 8.1 通常的容器构建方式

```Dockerfile
# First use this environment(Dockerfile.build) to compile our echo.go
FROM golang:1.16.3
WORKDIR /root/echosrv
COPY echo.go .
RUN CGO_ENABLED=0 go build -o echo echo.go
```

```Dockerfile
# This is our produnction environment
FROM alpine:latest
WORKDIR /root/
COPY echo .
CMD ["./echo"]
```

```Bash
#!/bin/sh
echo Building echo service
docker build --no-cache -t echo:build . -f Dockerfile.build     # build our binary in image(echo:build)
docker create --name echosrv echo:build                 # create new container(echosrv) for copying binary
docker cp echosrv:/root/echosrv/echo ./                 # copy our binary to current work directory
docker rm echosrv                                       # delete the temporary container 

echo Building echo:latest
docker build -t echo:latest .                           # since our binary is ready, we can build our echo server now.
```

构建并执行后测试发现运行正常。但采用这种方法需要维护两个Dockerfile以及一个build脚本，
其中应当有简化的空间，也就是下面要说的multi-stage构建。

```
➜  multi-stage docker run --detach=true --name=echo1 -p 5000:8000 echo:latest
0ec4d759a2b5318978d6fcfa9ca9e4d973c961de95c81d90b4e8a641dd993f36
➜  multi-stage curl "http://192.168.1.22:5000/echo?name=Rob&name=Ken&name=Robert"
name = Rob,Ken,Robert#
```

# 8.2 multi-stage构建

```Dockerfile
FROM golang:1.16.3
WORKDIR /root/echo_multi
COPY echo.go .
RUN CGO_ENABLED=0 go build -o echo echo.go

FROM alpine:latest
WORKDIR /root
COPY --from=0 /root/echo_multi/ .
CMD ["./echo"]
```

用--from=0引用第一个stage构建的镜像，但这样指定不是很清晰，于是命名式多阶段构建出现了。


# 8.3 named multi-stage构建

```Dockerfile
FROM golang:1.16.3 as builder
WORKDIR /root/echo_multi
COPY echo.go .
RUN CGO_ENABLED=0 go build -o echo echo.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/echo_multi/ .
CMD ["./echo"]
```

用FROM as来指定别名，然后用--from=来引用，除此之外，还有以下用法：

- docker build --target builder ...     // 在stage: builder处停止
- COPY --from=nginx:latest /etc/nginx/nginx.conf /nginx.conf        // 使用外部nginx:latest镜像作为stage

