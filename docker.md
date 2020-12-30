# Docker

## 总览

```
查询镜像（Images）、容器（Containers）和本地卷（Local Volumes）等空间使用大户的空间占用情况
docker system df

进一步查看空间占用详情
docker system df -v
```

通过 Docker 内置的 CLI 指令来进行自动空间清理:

```
docker system prune
```

该指令默认会清除所有如下资源：

- 已停止的容器（container）
- 未被任何容器所使用的卷（volume）
- 未被任何容器所关联的网络（network）
- 所有悬空镜像（image）。

增加`-a`参数后可一并清除所有未使用的镜像和悬空镜像,  `-f`参数用于忽略相关告警确认信息:

```
docker system prune -a -f
```



## 镜像

[官方镜像库](https://hub.docker.com/)

```
从官方搜索minio镜像
docker search minio

从官方下载minio镜像  可以不加:latest 默认自带的
docker pull minio/minio:latest

查看所有镜像
docker images

单个镜像删除，相当于：docker rmi redis:latest
docker rmi redis

强制删除(针对基于镜像有运行的容器进程)
docker rmi -f redis

多个镜像删除，不同镜像间以空格间隔
docker rmi -f redis tomcat nginx

删除全部镜像
docker rmi -f $(docker images -q)

删除悬空镜像
docker rmi $(docker images -q -f dangling=true)
或
docker image prune -f

删除所有未使用的镜像
docker image prune -f -a
```



## 卷

```
查看所有卷
docker volume ls

查看指定容器卷详情信息
docker volume inspect volume-mongodb

创建卷
docker volume create volume-mongodb

删除卷
docker volume rm volume-mongodb

清理无用卷
docker volume prune
```





## 容器

```
查看正在运行容器
docker ps

查看所有容器
docker ps -a

下面对容器的操作, 最后参数可以是容器name或者容器id 

停止一个运行中的容器
docker stop minio-chry

启动已经被停止的容器
docker start minio-chry

重启容器
docker restart minio-chry

杀掉一个运行中的容器
docker kill minio-chry

删除容器
docker rm minio-chry

进入容器后台终端, 需要提前知道该容器的终端类型是sh还是bash
docker exec -it minio-chry /bin/bash
```

### 容器的启动

[docker run 参数详解](https://blog.csdn.net/weixin_39998006/article/details/99680522)

[Docker run 命令参数及使用](https://blog.csdn.net/luolianxi/article/details/107169954)

```
命令格式：docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
```

常用选项:

```
-d, --detach=false， 指定容器运行于前台还是后台，默认为false
-i, --interactive=false， 打开STDIN，用于控制台交互
-t, --tty=false， 分配tty设备，该可以支持终端登录，默认为false
	-it 即表示可以命令行模式进入该容器 常用于linux系统容器
-p, --publish=[]， 指定容器暴露的端口 
	-p 80:8080 表示将宿主机的80端口映射到容器的8080端口
-v, --volume=[]， 给容器挂载存储卷，挂载到容器的某个目录
	-v /mnt/e/temp/minio/data:/data 表示将本地目录/mnt/e/temp/minio/data挂载到容器的/data目录
	-v volume-mongodb:/data/db 表示将本地卷volume-mongodb挂载到容器的/data/db目录

--name=""， 指定容器名字，后续可以通过名字进行容器管理，links特性需要使用名字
--restart="no"， 指定容器停止后的重启策略:
	no：容器退出时不重启
	on-failure：容器故障退出（返回值非零）时重启
		--restart="on-failure:10" 表示最多重启10次
	always：容器退出时总是重启
--rm=false， 指定容器停止后自动删除容器(有人说不支持以docker run -d启动的容器, 但本人实测确发现可以)
	docker: Conflicting options: --restart and --rm.


```

常用示例:

```
创建一个停止后自动删除容器
docker run -itd --rm=true personreidserver:0.1

创建一个停止后自动删除容器, 不执行默认命令直接进入其终端 (需要提前知道该容器的终端是sh还是bash, 而且这时候就不应该加入-d参数了)
docker run -it --rm=true personreidserver:0.1 /bin/sh

强烈建议所有创建都指定名称:
docker run -it --rm=true --name pserver personreidserver:0.1

```

## 网桥

[Docker容器间通信方法](https://juejin.cn/post/6844903847383547911)

```
创建了一个名为"my-net"的网络
docker network create my-net

将容器加入my-net网络中
docker network connect my-net test_demo  
docker network connect my-net mysqld5.7

查看my-net的网络配置
docker network inspect my-net

测试别名互连: 进入容器test_demo的终端, ping mysqld5.7

断开容器与docker0的连接
docker network disconnect bridge test_demo
docker network disconnect bridge mysqld5.7

```



# Dockerfile

Dockerfile 是一个用来构建镜像的文本文件，文本内容包含了一条条构建镜像所需的指令和说明。

## 语法

```dockerfile
# 思路: 分编译和发布两个阶段, 在带go语言的docker镜像中编译，将编译出来的二进制文件拷贝到一个不带go环境的较小的镜像
# 编译阶段过程中会产生名为none的镜像, 删除之:
#docker rmi $(docker images | grep "none" | awk '{print $3}')
#docker rmi $(docker images -q -f dangling=true)

# 运行示例: docker build -t personreidserver:0.1 -f Dockerfile ..
# -t personreidserver:0.1 指定构建的镜像名称和tag
# -f Dockerfile表示指定Dockerfile文件路径
# .. 表示上下文目录


#编译阶段：使用golang:rc-alpine 版本
FROM golang:rc-alpine as build

# 容器环境变量添加
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 设置当前工作区
WORKDIR /go/release

# 把全部文件添加到/go/release目录
# 第一个. 表示运行本Dockerfile时, docker build命令所指定的目录 一般要求是源码文件
# 第二个. 表示上面WORKDIR指定的目录
ADD . .

# [Alpine 的 CGO 问题](https://www.jianshu.com/p/22956db3a52a)
# [Docker与Golang的巧妙结合](http://dockone.io/article/1712)
# [在Goland中使用Docker插件生成镜像与创建容器](https://www.cnblogs.com/litchi99/p/13724811.html)
# [多阶段构建Golang程序Docker镜像](https://www.cnblogs.com/FireworksEasyCool/p/12838875.html)
# 编译: 把main.go编译为可执行的二进制文件, 并命名为app
# alpine只支持静态链接, 当cgo开启时，默认是按照动态库的方式来链接so文件, 于是使用CGO_ENABLED=0关闭了cgo
# 如果需要使用cgo, 可调用 go build --ldflags="-extldflags -static" 来让gcc使用静态编译解决问题
# GOARCH：32位系统为386，64位系统为amd64
# -ldflags参数: -w关闭所有告警信息 -s省略符号表和调试信息
# -installsuffix 在软件包安装的目录中增加后缀标识，用于区分默认版本, 比如这里指定为cgo
# -o：指定编译后的可执行文件名称
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -ldflags="-s -w" -installsuffix cgo -o app cmd/httpserver/main.go



# 上面将源码编译出了二进制文件, 下面把二进制文件和相关依赖打包整一个镜像


# 使用alpine作为基础镜像
FROM alpine as production

# 拷贝本机输入目录下的release配置文件和资源文件
ADD ./release /release/
# 复制build阶段编译出来的可执行二进制文件到production阶段的/release/目录下
COPY --from=build /go/release/app /release/

# 指定运行目录
WORKDIR /release
# 启动服务
CMD ["./app"]

```

## 运行

```
docker build -t name:tag -f filepath dir
-t指定构建的镜像名和TAG
-f指定构建使用的Dockerfile路径
最后一个参数指定上下文目录
如:
docker build -t pdserver:0.1 -f docker/Dockerfile ./source_code
```





# Docker Compose





# NVIDIA Docker

注意: 不支持在虚拟机上安装NVIDIA显卡驱动

[Ubuntu16.04下安装NVIDIA显卡驱动](https://blog.csdn.net/yinwangde/article/details/89439648)

[可能会遇到的问题](https://blog.csdn.net/weixin_43002433/article/details/108888927)

[Docker使用篇之Nvidia-docker](https://blog.csdn.net/felaim/article/details/105229226)

## 运行gpu容器

```
使用所有GPU
docker run --gpus all nvidia/cuda:9.0-base nvidia-smi

使用两个GPU
docker run --gpus 2 nvidia/cuda:9.0-base nvidia-smi

指定GPU运行
docker run --gpus '"device=1,2"' nvidia/cuda:9.0-base nvidia-smi
docker run --gpus '"device=UUID-ABCDEF,1"' nvidia/cuda:9.0-base nvidia-smi
```

## 在docker-compose中使用gpu

```
version: '2.4'
services:
  nvsmi:
    image: ubuntu:16.04
    runtime: nvidia
    environment:
      - NVIDIA_VISIBLE_DEVICES=all
    command: nvidia-smi
```



