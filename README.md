# goweb
## 代码分层 CLD
```
Controller  对外接口 转换协议
----------------------------
Logic 业务
----------------------------
DAO  data access object
```
# docker
```
从官方搜索minio镜像
docker search minio

从官方下载三个数据库的镜像  可以不加:latest 默认自带的
docker pull minio/minio:latest
docker pull redis:latest
docker pull mongo:latest

查看所有镜像
docker images

##单个镜像删除，相当于：docker rmi redis:latest
docker rmi redis
##强制删除(针对基于镜像有运行的容器进程)
docker rmi -f redis
##多个镜像删除，不同镜像间以空格间隔
docker rmi -f redis tomcat nginx
##删除本地全部镜像
docker rmi -f $(docker images -q)
------------------------------------------------------------------------------------
##启动容器
-i: 交互式操作。
-t: 临时终端。
-it即表示启动后以命令行模式进入该容器 常用于linux系统容器
-d: 表示后台运行 
-p: 表示暴露端口,映射容器服务的9000端口到宿主机的9000端口。外部可以直接通过宿主机ip:9000访问到minio的服务。
-v: 给容器挂载存储卷，挂载到容器的某个目录    
尽量用--name指定进程名称
然后注意镜像名称如minio/minio,后面跟上的都是在镜像系统里面运行的命令,相当于在linux终端上运行"server /data"
docker run -p 9000:9000 --name minio-chry -d \
-v /mnt/e/temp/minio/data:/data \
minio/minio server /data

##查看所有容器 -a表示显示已经停止的容器
chry@DESKTOP-N2FV3PF:/mnt/c/Users/a8512$ docker ps -a
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                       PORTS               NAMES
ac6e62eb0d34        minio/minio         "/usr/bin/docker-ent…"   5 minutes ago       Exited (137) 3 seconds ago                       redis-chry
17f69a7ae3c4        minio/minio         "/usr/bin/docker-ent…"   30 minutes ago      Exited (0) 6 minutes ago                         elegant_dijkstra

##停止一个运行中的容器    最后参数是--name指定的名称或者docker ps -a查看
docker stop minio-chry
##启动已经被停止的容器 最后参数是--name指定的名称或者docker ps -a查看
docker start minio-chry
##根据elegant_dijkstra的CONTAINER ID删除容器  
docker rm 17f69a7ae3c4
##重启容器
docker restart minio-chry
##杀掉一个运行中的容器
docker kill minio-chry

```



# redis
[Redis配置数据持久化](https://blog.csdn.net/ljl890705/article/details/51039015)
[docker安装并运行redis](https://blog.csdn.net/weixin_38424794/article/details/104301969)
```
安装
sudo apt install redis


运行服务器
redis-server
或者
docker run -p 6379:6379 --name redis-chry -d \
-v /mnt/e/temp/redis/data:/data \
-v /mnt/e/temp/redis/config:/etc/redis \
redis redis-server --appendonly yes
或者 不指定配置文件
docker run -p 6379:6379 --name redis-chry -d \
-v /mnt/e/temp/redis/data:/data \
redis redis-server --requirepass "chry" --appendonly yes



进入客户端
redis-cli
或者
chry@DESKTOP-N2FV3PF:/mnt/c/Users/a8512$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
c243fe394927        redis               "docker-entrypoint.s…"   21 seconds ago      Up 19 seconds       0.0.0.0:6379->6379/tcp   redis-chry
33ecaecdf556        minio/minio         "/usr/bin/docker-ent…"   8 minutes ago       Up 8 minutes        0.0.0.0:9000->9000/tcp   minio-chry
chry@DESKTOP-N2FV3PF:/mnt/c/Users/a8512$ docker exec -it  redis-chry redis-cli
```

以下操作在redis客户端shell操作
```
查看密码
127.0.0.1:6379> config get requirepass
1) "requirepass"
2) ""

设置密码
127.0.0.1:6379> config set requirepass chry
OK

后面的操作就要先验证密码, 否则会提示NOAUTH
127.0.0.1:6379> config get requirepass
(error) NOAUTH Authentication required.

验证密码
127.0.0.1:6379> AUTH chry
OK
127.0.0.1:6379> config get requirepass
1) "requirepass"
2) "chry"
```

常用命令
```
设置key-val
SET resource:lock "Redis Demo"

设置key生存时间(秒)
EXPIRE resource:lock 120

查看key剩余时间
TTL resource:lock

设置key-val和生存时间(秒)
SET resource:lock "Redis Demo 3" EX 5

删除过期值并使该键再次永久存在
PERSIST resource:lock

```

# minio
```
docker启动minio, 默认账密adminminio/adminminio
docker run -p 9000:9000 --name minio-chry -d -v /mnt/e/temp/minio/data:/data minio/minio server /data
```

# mongo
```
docker for windows使用mongodb镜像，如果直接使用 -v 参数挂载磁盘目录，启动镜像的时候会报错。
原因:Windows和OS X上的默认Docker设置使用VirtualBox VM来托管Docker守护程序。不幸的是，VirtualBox用于在主机系统和Docker容器之间共享文件夹的机制与MongoDB使用的内存映射文件不兼容（请参阅vbox bug，docs.mongodb.org和相关的jira.mongodb.org错误）。这意味着无法运行映射到主机的数据目录的MongoDB容器。
解决:
1，使用docker命令创建卷： docker volume create volume-mongodb
2，然后挂载到上一步创建的卷：
docker run -p 27017:27017 --name mongo-chry -d -v volume-mongodb:/data/db mongo
```