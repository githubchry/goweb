

# docker环境

参考链接:

[Redis配置数据持久化](https://blog.csdn.net/ljl890705/article/details/51039015)

[docker安装并运行redis](https://blog.csdn.net/weixin_38424794/article/details/104301969)

```shell
创建 redis.conf 下载http://download.redis.io/redis-stable/redis.conf
文件放到/mnt/e/docker/mnt/redis/config/redis.conf

docker run -p 6379:6379 --name redis-chry -d \
-v /mnt/e/docker/mnt/redis/data:/data \
-v /mnt/e/docker/mnt/redis/config:/etc/redis \
redis redis-server --appendonly yes --requirepass "chry"

--appendonly yes		表示开启数据持久化
--requirepass "chry"	表示以"chry"作为授权密码
```

## 可用性测试
进入redis shell

```shell
docker exec -it redis-chry redis-cli
```

以下操作在redis客户端shell操作

```
如果以--requirepass方式启动容器, 进入redis shell后需要先验证密码, 否则会提示NOAUTH
127.0.0.1:6379> config get requirepass
(error) NOAUTH Authentication required.

验证密码
127.0.0.1:6379> AUTH chry
OK

查看密码
127.0.0.1:6379> config get requirepass
1) "requirepass"
2) "chry"

修改密码
127.0.0.1:6379> config set requirepass chry
OK

查看密码
127.0.0.1:6379> config get requirepass
1) "requirepass"
2) ""

如果重新进入redis shell, 后面的操作就要先验证密码, 否则会提示NOAUTH
127.0.0.1:6379> config get requirepass
(error) NOAUTH Authentication required.

```

## docker-compose配置

```
# 版本3支持集群部署 版本2仅支持单机部署
version: '3'
services:
  redis:
    image: redis:latest
    container_name: redis
    restart: always
    ports:
      - "6379:6379"
    volumes:
      - /mnt/e/docker/mnt/redis/data:/data
      - /mnt/e/docker/mnt/redis/config:/etc/redis
    #指定授权密码, 持久化数据
    command: redis-server --appendonly yes --requirepass "chry"
```



# 常用命令

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



# 异常断电

```
进入容器或者新建容器的终端，执行恢复操作

docker run -it --rm=true -v /app/data/redis/data:/data redis /bin/sh

cd /data
redis-check-aof --fix appendonly.aof
```

