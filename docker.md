# docker

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



# dockerfile









# docker-compose

