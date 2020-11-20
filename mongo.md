

# docker环境

参考链接:

[使用Docker创建MongoDb服务](https://www.cnblogs.com/bowendown/p/12656380.html)

[使用Docker一键部署MongoDB](https://blog.csdn.net/u011104991/article/details/81735960)

## win10下挂载db目录问题

```
docker for windows使用mongodb镜像，如果直接使用 -v 参数挂载磁盘目录，启动镜像的时候会报错。
原因:Windows和OS X上的默认Docker设置使用VirtualBox VM来托管Docker守护程序。不幸的是，VirtualBox用于在主机系统和Docker容器之间共享文件夹的机制与MongoDB使用的内存映射文件不兼容（请参阅vbox bug，docs.mongodb.org和相关的jira.mongodb.org错误）。这意味着无法运行映射到主机的数据目录的MongoDB容器。
解决:
1，使用docker命令创建卷： 
docker volume create volume-mongodb
2，然后挂载到上一步创建的卷：
docker run -p 27017:27017 --name mongo-chry -d -v volume-mongodb:/data/db mongo
```

## 关于--auth参数

`mongodb`默认不启用授权认证，只要能连接到该服务器，就可连接到`mongod`, 若要启用安全认证，需要增加`--auth`参数.

以`--auth`方式启动`mongodb`后, 需要使用用户名和密码进行授权才能操作数据库, 但如果在`mongodb`还未创建任何用户的情况下, 是无法进行auth的, 也就无法进行任何操作. 

网上常用的解决方法:

1. 以默认方式创建一个mongo容器, 同时挂载db目录, 然后登录db创建用户, 最后退出并销毁容器

   ```shell
   主机终端下使用默认配置创建并启动mongo容器,同时挂载db目录:
   docker run -p 27017:27017 --name mongo-chry -v volume-mongodb:/data/db -d mongo
   主机终端下进入到mongo-chry容器里面的mongodb终端:
   docker exec -it mongo-chry mongo
   创建超级用户
   use admin
   db.createUser({user:"admin", pwd:"admin",roles:["root"]})
   
   (在chrydb下新建chry用户并授予所有者权限)
   use chrydb
   db.createUser({user:"chry", pwd:"chry",roles:[{role:"dbOwner", db:"chrydb"}]})
   
   exit
   执行exit后已经退出到主机, 停止并销毁容器
   docker stop mongo-chry && docker rm mongo-chry
   ```

2. 此时上一步创建的用户信息依旧会保存在db目录, 以--auth方式创建按mongo容器, 同时挂载db目录即可

   ```shell
   docker run -p 27017:27017 --name mongo-chry -v volume-mongodb:/data/db -d mongo --auth
   
   主机终端下进入到mongo-chry容器里面的mongodb终端测试一下:
   docker exec -it mongo-chry mongo
   
   > use chrydb
   > show dbs
   啥都看不到, 因为auth模式下先验证用户才能进行其他操作
   > db.auth('chry','chry')
   1
   > show dbs
   chrydb  0.000GB
   > exit
   bye
   ```

以上随为常规操作,  但仍很麻烦, 无法达到一键部署的目的,  本帅肯定不能忍, 所以找到了曲线救国的方法......

## 曲线救国解决auth模式一键部署问题

先把旧的干掉:

```
docker stop mongo-chry && docker rm mongo-chry && docker volume rm volume-mongodb
```

mongo的镜像系统有这么一个机制, 当容器首次启动时它会以脚本字母顺序执行`/docker-entrypoint-initdb.d`目录下的`sh和js脚本`.

同时`js脚本`将被使用`MONGO_INITDB_DATABASE`指定的库或默认自带的test库执行 。

以上, 首先创建一个`setup.js`文件, 用来初始化创建我们的基础db信息:

```js
// 创建一个名为"chrydb"的DB
db = db.getSiblingDB('chrydb');

// 创建一个名为"chry"的用户，设置密码和权限
db.createUser({user: "chry", pwd: "chry", roles: [{ role: "dbOwner", db: "chrydb"}]});

// 在"chry"中创建一个名为"chry"的Collection 因为至少创建一个集合才能auth后执行show dbs看到 可省略
db.createCollection("chry");
```

把该文件放在目录下比如`/mnt/e/ubuntu/codes/git/goweb/release/mongo`, 在创建容器时挂载到`/docker-entrypoint-initdb.d`目录:

```
docker run -p 27017:27017 --name mongo-chry -d -v /mnt/e/ubuntu/codes/git/goweb/release/mongo:/docker-entrypoint-initdb.d/ mongo --auth
```

测试一下:

```shell
docker exec -it mongo-chry mongo chrydb
进入后台后执行:
> show dbs
> db.auth('chry','chry')
1
> show dbs
chrydb  0.000GB
> exit
bye
```

鼓掌.

## docker-compose配置

```
# 版本3支持集群部署 版本2仅支持单机部署
version: '3'
#声明volumes
volumes:
  volume-mongodb: { }
services:
  mongo:
    image: mongo:latest
    container_name: mongo
    restart: always
    ports:
      - "27017:27017"
    volumes:
      # 由于兼容问题无法直接挂载win目录到mongo, 只能创建一个volume后挂载
      #- volume-mongodb:/data/db
      - ./mongo:/docker-entrypoint-initdb.d/
    # 启动授权登录
    command: --auth
```

## 拓展: 环境变量

[docker 官方mongodb镜像](https://www.jianshu.com/p/9703e77e7931)

```

```





# 基础操作

## 进入mongo shell

```
终端下直接运行:
mongo

进入shell后切换到chrydb => mongo db名
例:
mongo chrydb

进入容器里面的mongo shell => docker exec -it 容器名 mongo [db名]
例:
docker exec -it mongo-chry mongo chrydb
```

## db - user - collection - document

```


创建超级用户
use admin
db.createUser({user:"admin", pwd:"admin",roles:[{role: "root", db: "admin"}]})

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

创建collection(集合)
> db.createCollection("trainer")
{ "ok" : 1 }
> db.createCollection("student")
{ "ok" : 1 }

显示所有collection
>  show collections
student
trainer

删除指定数据集
> db.trainer.drop()
true

插入一条文档
> db.student.insertOne({name:"小王子",age:18});
{
        "acknowledged" : true,
        "insertedId" : ObjectId("5f5848ee38ad7149539434e9")
}

插入多条文档
> db.student.insertMany([
... {name:"张三",age:20},
... {name:"李四",age:25}
... ]);
{
	"acknowledged" : true,
	"insertedIds" : [
		ObjectId("5f58491d38ad7149539434ea"),
		ObjectId("5f58491d38ad7149539434eb")
	]
}

查询所有文档：
> db.student.find()
{ "_id" : ObjectId("5f5848ee38ad7149539434e9"), "name" : "小王子", "age" : 18 }
{ "_id" : ObjectId("5f58491d38ad7149539434ea"), "name" : "张三", "age" : 20 }
{ "_id" : ObjectId("5f58491d38ad7149539434eb"), "name" : "李四", "age" : 25 }

查询age>20岁的文档：
> db.student.find(
... {age:{$gt:20}}
... )
{ "_id" : ObjectId("5f58491d38ad7149539434eb"), "name" : "李四", "age" : 25 }

更新文档：
> db.student.update({name:"小王子"},{name:"老王子",age:98})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })

删除文档：
db.student.deleteOne({name:"李四"});
```



