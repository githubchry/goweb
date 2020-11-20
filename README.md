# goweb
运行时确保工作目录为:goweb/cmd

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
创建 redis.conf 下载http://download.redis.io/redis-stable/redis.conf
文件放到这里/mnt/e/temp/redis/config/redis.conf 
找到bind 127.0.0.1，把这行前面加个#注释掉

docker run -p 6379:6379 --name redis-chry -d \
-v /mnt/e/temp/redis/data:/data \
-v /mnt/e/temp/redis/config:/etc/redis \
redis redis-server --appendonly yes
或者 不指定配置文件
docker run -p 6379:6379 --name redis-chry -d \
-v /mnt/e/temp/redis/data:/data \
redis redis-server --appendonly yes --requirepass "chry"



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
docker run -p 27017:27017 --name mongo-chry -d -v volume-mongodb:/data/db mongo --auth

进入mongo shell
docker exec -it mongo-chry mongo

创建超级用户
use admin
db.createUser({user:"admin", pwd:"admin",roles:["root"]})
 
查看dbs列表
> show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
test    0.000GB

切换到其中一个db
> use test
switched to db test
为当前db创建用户
db.createUser({user:"chry", pwd:"chry",roles:["readWrite"]})


后面再进入就要验证用户
db.auth('chry','chry')


```


# docker-compose
[基于Docker + Go+ Kafka + Redis + MySQL的秒杀已经Jmeter压力测试](https://blog.csdn.net/q3585914/article/details/90604565)

[kafka的Docker镜像使用说明(wurstmeister/kafka)](https://blog.csdn.net/boling_cavalry/article/details/85395080)

[《KAFKA官方文档》入门指南](http://ifeve.com/kafka-1/)

[kafka跨网段和外网访问](https://segmentfault.com/a/1190000020715650)

Docker Compose是一个用来定义和运行复杂应用的Docker工具。一个使用Docker容器的应用，通常由多个容器组成。使用Docker Compose不再需要使用shell脚本来启动容器。 
Compose 通过一个配置文件来管理多个Docker容器，在配置文件中，所有的容器通过services来定义，然后使用docker-compose脚本来启动，停止和重启应用，和应用中的服务以及所有依赖服务的容器，非常适合组合使用多个容器进行开发的场景。

比如, kafka依赖于zookeeper, 每次使用kafka前都要先部署zookeeper
```shell script
docker run -d --name zookeeper --publish 2181:2181  wurstmeister/zookeeper
docker run -d --name kafka --publish 9092:9092 \
--link zookeeper \
--env KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
--env KAFKA_ADVERTISED_HOST_NAME=192.168.50.16 \
--env KAFKA_ADVERTISED_PORT=9092 \
wurstmeister/kafka


// 注意上面的KAFKA_ADVERTISED_HOST_NAME, 在跨IP的情况下不要填127.0.0.1
```
于是可以通过docker-compose直接把这两个容器管理起来:

```yaml
version: '3'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - "2181:2181"
  kafka:
    image: wurstmeister/kafka
    build: .
    ports:
      - "9092:9092"
    environment:
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://10.11.5.90:9092
      KAFKA_LISTENERS: PLAINTEXT://:9092
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    volumes:
      - /mnt/e/temp/kafka/docker.sock:/var/run/docker.sock
```
在`docker-compose.yml`当前目录下运行:`docker-compose up -d`
```shell script
chry@DESKTOP-N2FV3PF:/mnt/e/temp/kafka$ docker-compose up -d
Creating network "kafka_default" with the default driver
Creating kafka_kafka_1     ... done
Creating kafka_zookeeper_1 ... done
chry@DESKTOP-N2FV3PF:/mnt/e/temp/kafka$

进入kafka_kafka_1后台测试可用性:
docker exec -it kafka_kafka_1 bash

## 创建主题
kafka-topics.sh --create --zookeeper 10.11.5.89:2181 --replication-factor 1 --partitions 1 --topic event1
## 查看主题 
kafka-topics.sh --list --zookeeper 10.11.5.100:2181
## 发送消息
kafka-console-producer.sh --broker-list 10.11.5.100:9092 --topic mykafka
## 接受消息
kafka-console-consumer.sh --bootstrap-server 10.11.5.100:9092 --from-beginning --topic event1

## 删除主题 
kafka-topics.sh --delete --zookeeper 10.11.5.100:2181  --topic mykafka

退出并删除
docker-compose down
```


[证书访问](https://www.jianshu.com/p/3102418e5a7d)
#


