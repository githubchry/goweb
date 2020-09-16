# goweb

# redis

```
安装
sudo apt install redis

运行服务器
redis-server

进入客户端
redis-cli
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
docker run -p 9000:9000 minio/minio server /data
```
