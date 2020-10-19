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
/mnt/e/temp/minio/algorithm:/mnt

docker run  -itd --gpus all -p 8080:8080 72439a283e79
docker run  -itd -p 8080:8080 --name algorithm-cpu -d -v /mnt/e/temp/algorithm/mnt:/mnt 72439a283e79

进入后台
docker exec -it algorithm-cpu /bin/bash

export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPATH=/mnt/gopath

cd /mnt/goweb/internal/logics/algorithm/person_detection
gcc -c person_detection.cpp person_detection_wrapper.cpp algorithm.pb.cc -I/mnt/new-root/usr/local/include
