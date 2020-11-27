

[TOC]



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
docker run -p 17017:27017 --name mongo-chry -d -v /mnt/e/ubuntu/codes/git/goweb/release/mongo:/docker-entrypoint-initdb.d/ mongo --auth
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



# mongo shell

[The mongo Shell](https://docs.mongodb.com/manual/mongo/)

```
默认情况下连接到本机mongo shell, 直接在终端运行:
mongo

连接远端的mongo shell => mongo 地址:端口号
mongo 192.168.58.2:27017

进入shell后切换到chrydb => mongo db名
mongo chrydb

指定库，用户和密码连接mongo shell
mongo chrydb -uchry -pchry

进入容器里面的mongo shell => docker exec -it 容器名 mongo [参数]
例:
docker exec -it mongo-chry mongo chrydb -uchry -pchry
```

## 

# 语法命令

```
所有MongoDB的组合单词都使用首字母小写的驼峰式写法
db.createCollection("xxx")


使用类似javascript语法，对象会自动补全为规范的json格式:
{name:"张龙豪"，age:18}    =>   {"name":"张龙豪","age":18}


```

## 数据类型

| **Type**                | **Number** |
| ----------------------- | ---------- |
| Double                  | 1          |
| String                  | 2          |
| Object                  | 3          |
| Array                   | 4          |
| Binary data             | 5          |
| Object id               | 7          |
| Boolean                 | 8          |
| Date                    | 9          |
| Null                    | 10         |
| Regular Expression      | 11         |
| JavaScript              | 13         |
| Symbol                  | 14         |
| JavaScript (with scope) | 15         |
| 32-bit integer          | 16         |
| Timestamp               | 17         |
| 64-bit integer          | 18         |
| Min key                 | 255        |
| Max key                 | 127        |

两个不同类型的值相比较时，按照如下顺序决定大小

1. MinKey (internal type)
2. Null
3. Numbers (ints, longs, doubles)
4. Symbol, String
5. Object
6. Array
7. BinData
8. ObjectID
9. Boolean
10. Date, Timestamp
11. Regular Expression
12. MaxKey (internal type)

当使用$type判断某个文档属性是否是MinKey时，不应使用255，应使用-1



## 基本操作

### 可用性测试

```
查看当前db
> db
test
>

查看所有数据库
> show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
以上三个都是系统自带的db, 与及默认登录时的test db(因为是空的没有集合所以不会显示在show dbs列表下)
config用处未知...
admin主要存储用户、角色等信息, 用户需要创建数据库帐号，访问时根据帐号信息来鉴权，而数据库帐号信息就存储在admin数据库下。
local主要存储副本集的配置信息、oplog信息，这些信息是每个Mongod进程独有的，不需要同步到副本集种其他节点。在使用MongoDB时，重要的数据千万不要存储在local数据库中，否则当一个节点故障时，存储在local里的数据就会丢失。

平时我们应该有意识地避免使用以上关键字作为数据库!
平时我们应该有意识地避免使用以上关键字作为数据库!
平时我们应该有意识地避免使用以上关键字作为数据库!
```

### 用户角色

[MongoDB 4.X 用户操作及角色](https://www.cnblogs.com/zhm1985/articles/13441133.html)

[MongoDB配置用户账号与访问控制](https://blog.csdn.net/qq_33206732/article/details/79877948)

#### 创建超级用户

```
use admin
db.createUser({user:"admin", pwd:"admin",roles:[{role: "root", db: "admin"}]})
```

## db - user - collection - document

#### 数据库db

```
切换到chrydb数据库(不存在则会自动创建, 但创建后不会显示在show dbs的结果里, 直到在该db下创建了集合)
> use chrydb
switched to db chrydb

删除当前数据库(比较少用)
db.dropDatabase()
```

#### 用户user 

```
use chrydb
为chrydb创建所有者用户:
db.createUser({user: "chry", pwd: "chry", roles: [{ role: "dbOwner", db: "chrydb"}]});
如果是给程序做普通的增删改查操作, 权限不宜太高, 就读写即可:
db.createUser({user: "chry", pwd: "chry" roles: ["readWrite"]});
```

#### 集合collection 

```
创建集合
> db.createCollection("student")
{ "ok" : 1 }
显示当前db下的所有集合
> show collections
student
删除student集合
> db.student.drop()
true
```

#### 文档document

```
插入一条文档(如果集合student不存在会自动创建)
> db.student.insertOne({name:"小王子",age:18})
{
        "acknowledged" : true,
        "insertedId" : ObjectId("5f5848ee38ad7149539434e9")
}

插入多条文档
> db.student.insertMany([
... {name:"张三",age:20},
... {name:"李四",age:25},
... {name:"李四",age:15}
... ])
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
{ "_id" : ObjectId("5f58491d38ad7149539434ec"), "name" : "李四", "age" : 15 }

结果格式化输出(语句后直接加.pretty()):
db.student.find().pretty()

查询age>20岁的文档：
> db.student.find({age:{$gt:20}})
{ "_id" : ObjectId("5f58491d38ad7149539434eb"), "name" : "李四", "age" : 25 }

查看年龄最大的人
db.student.find().sort({"age" : -1}).limit(1)
查看年龄最小的人
db.student.find().sort({"age" : 1}).limit(1)
查看特定地段不同值列表(同名的人当成一个)
db.student.distinct("name")

更新文档："score" : [ { "math" : 90 }, { "art" : 80 } ]
> db.student.update({name:"小王子"},{name:"老王子", age:98, score:{math: 90, art:80}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })

查询嵌套对象的值
> db.student.find({"score.math":90})
> db.student.find({"score.math":{$gt:60}})
注意这里的"score.math"一定要加双引号


删除文档：
db.student.deleteOne({name:"老王子"});
```

### 更新操作符

上面更新小王子对应的文档时,  用一个新的文档替换之, 如果只想修改原文档中的某一项, 怎么办? 

又或者想把原文档中的某一项删除, 怎么办? 

如果想把原文档中的age的值进行+1操作呢?

更新操作符可以帮助我们解决这些问题, 测试前先弄一条数据出来:

```
db.student.deleteOne({name:"chenzhou"});
db.student.insertOne({name:"chenzhou", age:22, birthday: "19941208"})
```



#### $inc 对数字字段做加法运算

```
用法：{$inc:{field:value}}

作用：对一个数字字段的某个field增加value

示例：将name为chenzhou的学生的age增加5

> db.student.update({name:"chenzhou"},{$inc:{age: 5}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb7f2bb1fb44be3a7a0412b"), "name" : "chenzhou", "age" : 27, "birthday" : "19941208" }
>
```

#### $min/$max

```
用法：{$min:{field:value,...}}    {$max:{field:value,...}}

作用：$max在field的值小于value时更新,  $min在field的值大于value时更新

示例：更新name为chenzhou的学生的age, 如果

> db.student.update({name:"chenzhou"},{$max:{age:30}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 30, "birthday" : "19941208" }
>
>
>
> db.student.update({name:"chenzhou"},{$min:{age:20}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 20, "birthday" : "19941208" }
```



#### $set 更新某字段

```
用法：{$set:{field:value}}

作用：把文档中某个字段field的值设为value

示例： 把chenzhou的年龄设为字符串表示的"23"

> db.student.update({name:"chenzhou"},{$set:{age:"23"}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb7f2bb1fb44be3a7a0412b"), "name" : "chenzhou", "age" : "23", "birthday" : "19941208" }
>
```



#### $unset 删除某字段

```
用法：{$unset:{field:1}}

作用：删除某个字段field

示例： 将chenzhou的年龄字段删除

> db.student.update({name:"chenzhou"},{$unset:{age:1}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb7f2bb1fb44be3a7a0412b"), "name" : "chenzhou", "birthday" : "19941208" }
>
```



#### $rename修改字段名

```
用法：{$rename:{old_field:"new_field"}}

作用：对字段进行重命名, 注意新字段名要加双引号, 旧字段名可加可不加

示例：把chenzhou记录的birthday字段重命名为bday

> db.student.update({name:"chenzhou"},{$rename:{birthday:"bday"}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb7f2bb1fb44be3a7a0412b"), "name" : "chenzhou", "bday" : "19941208" }
>
```

#### $currentDate置为当前时间

```
用法：{$currentDate:{field:true}} 或 {$currentDate:{field:{$type:"timestamp"}}}

作用：将字段的值置为当前日期/时间戳

示例：把chenzhou记录的bday置为当前日期/时间戳

> db.student.update({name:"chenzhou"},{$currentDate:{bday:true}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb7f2bb1fb44be3a7a0412b"), "name" : "chenzhou", "birthday" : "19941208", "bday" : ISODate("2020-11-20T17:00:42.673Z") }
> 
> 
> db.student.update({name:"chenzhou"},{$currentDate:{bday:{$type:"timestamp"}}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb7f2bb1fb44be3a7a0412b"), "name" : "chenzhou", "birthday" : "19941208", "bday" : Timestamp(1605922759, 1) }
>

```



#### $push追加元素到数组 - 与$slice/$sort/$each一起使用

[Mongo更新数组$slice修饰符](https://blog.csdn.net/yaomingyang/article/details/78698211)

```
用法：{$push:{field:value}}

作用：把value追加到field里。注：field只能是数组类型，如果field不存在，会自动插入一个数组类型

示例：给chenzhou添加别名"michael"

> db.student.update({name:"chenzhou"},{$push:{"ailas":"Michael"}})
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 20, "birthday" : "19941208", "ailas" : [ "Michael" ] }

```

##### $push+$each

```
如果需要一次性push多个元素, 可以使用$each配合
用法：{$push:{field:{"$each":value_array}}}

作用：把value_array追加到field里。注：field只能是数组类型，如果field不存在，会自动插入一个数组类型

示例：给chenzhou添加别名"A1", "A2", 添加成绩表{math:90}, {art:80}

> db.student.update({name:"chenzhou"},{$push:{"ailas":{$each:["A1", "A2"]}}})
> db.student.update({name:"chenzhou"},{$push:{"score":{$each:[{math:90}, {art:80}]}}})

> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 20, "birthday" : "19941208", "ailas" : [ "Michael", "A1", "A2" ], "score" : [ { "math" : 90 }, { "art" : 80 } ] }
>
```

##### $push+$each+$slice

插入数组后截取: `$slice: -5`表示截取后面5个; `$slice: 5`表示截取前面5个;

```
一个students集合文档如下：
{ "_id" : 1, "scores" : [ 40, 50, 60 ] }
使用修饰符对集合文档进行操作：
db.students.update(
   { _id: 1 },
   {
     $push: {
       scores: {
         $each: [ 80, 78, 86 ],
         $slice: -5
       }
     }
   }
)
操作的结果是取得了后面五个元素：
{ "_id" : 1, "scores" : [  50,  60,  80,  78,  86 ] }

```

如果$each数组为空则表示不插入截取

```
一个students集合文档如下：
{ "_id" : 3, "scores" : [  89,  70,  100,  20 ] }
如果$each为空的情况下:
db.students.update(
  { _id: 3 },
  {
    $push: {
      scores: {
         $each: [ ],
         $slice: -3
      }
    }
  }
)
结果是截取了后三个元素
{ "_id" : 3, "scores" : [  70,  100,  20 ] }
```

##### $push+$each+$slice+$sort

插入数组后排序再截取, `$sort: { score: -1 }`表示倒序, `$sort: { score: 1 }`表示正序, 

```
一个students集合文档如下：
{
   "_id" : 5,
   "quizzes" : [
      { "wk": 1, "score" : 10 },
      { "wk": 2, "score" : 8 },
      { "wk": 3, "score" : 5 },
      { "wk": 4, "score" : 6 }
   ]
}
使用修饰符对集合文档进行操作：
db.students.update(
   { _id: 5 },
   {
     $push: {
       quizzes: {
          $each: [ { wk: 5, score: 8 }, { wk: 6, score: 7 }, { wk: 7, score: 6 } ],
          $sort: { score: -1 },
          $slice: 3
       }
     }
   }
)
操作结果如下：
{
  "_id" : 5,
  "quizzes" : [
     { "wk" : 1, "score" : 10 },
     { "wk" : 2, "score" : 8 },
     { "wk" : 5, "score" : 8 }
  ]
}
```

#### $addToSet值不存在时才追加到数组

```
用法：{$addToSet:{field:value}}

作用：加一个值到数组内，而且只有当这个值在数组中不存在时才增加。

示例：往chenzhou的别名字段里添加别名A2, A3

> db.student.update({name:"chenzhou"},{$addToSet:{"ailas":"A2"}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 0 })
> db.student.update({name:"chenzhou"},{$addToSet:{"ailas":"A3"}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 20, "birthday" : "19941208", "ailas" : [ "Michael", "A1", "A2", "A3" ], "score" : [ { "math" : 90 }, { "art" : 80 } ] }
>



分两次不够优雅, 类似$push, $addToSet也可以与$each组合:
> db.student.update({name:"chenzhou"},{$addToSet:{"ailas":{$each:["A2", "A3"]}}})

```

#### $pop删除数组第一个或最后一个值

```
用法：删除数组内第一个值：{$pop:{field:-1}}、删除数组内最后一个值：{$pop:{field:1}}

作用：用于删除数组内的一个值

示例： 删除chenzhou记录中alias字段中第一个和最后一个别名

> db.student.update({name:"chenzhou"},{$pop:{"ailas":-1}})
> db.student.update({name:"chenzhou"},{$pop:{"ailas":1}})

> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 20, "birthday" : "19941208", "ailas" : [ "A1" ], "score" : [ { "math" : 90 }, { "art" : 80 } ] }
>
```



#### $pull删除数组中所有等于value的值

```
用法：{$pull:{field:_value}}

作用：从数组field内删除所有等于_value的值

先造点数据
> db.student.update({name:"chenzhou"},{$push:{"ailas":{$each:["A1", "A2", "A3"]}}})
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 20, "birthday" : "19941208", "ailas" : [ "A1", "A1", "A2", "A3" ], "score" : [ { "math" : 90 }, { "art" : 80 } ] }


示例：删除chenzhou记录中的别名A1

> db.student.update({name:"chenzhou"},{$pull:{"ailas":"A1"}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.student.find({name:"chenzhou"})
{ "_id" : ObjectId("5fb871351fb44be3a7a0412c"), "name" : "chenzhou", "age" : 20, "birthday" : "19941208", "ailas" : [ "A2", "A3" ], "score" : [ { "math" : 90 }, { "art" : 80 } ] }
>

```

https://blog.csdn.net/yaomingyang/article/details/78698211

### 条件操作符

updata更新和find查找时可以通过条件操作符指定过滤条件

以下为常用的表达式：

#### 比较表达式

| 操作符 | 语义                         | 语法                   | 示例                               |
| :----- | :--------------------------- | ---------------------- | ---------------------------------- |
|        | equals等于 ==                | {<key>:<value>}        | db.student.find({name:"chenzhou"}) |
| $ne    | not equals不等于!=           | {<key>:{$ne:<value>}}  | db.student.find({age:{$ne:50}})    |
| $gt    | great than大于 >             | {<key>:{$gt:<value>}}  | db.student.find({age:{$gt:50}})    |
| $gte   | great than equals大于等于 >= | {<key>:{$gte:<value>}} | db.student.find({age:{$gte:50}})   |
| $lt    | less than小于 <              | {<key>:{$lt:<value>}}  | db.student.find({age:{$lt:50}})    |
| $lte   | less than equals小于等于 <=  | {<key>:{$lte:<value>}} | db.student.find({age:{$lte:50}})   |

#### 逻辑表达式

| 操作符 | 语义       | 语法                                              | 示例                                      |
| :----- | :--------- | ------------------------------------------------- | ----------------------------------------- |
|        | 隐式逻辑与 | {<key1>:<value1>, <key2>:<value2>, ...}           | db.xx.find({name:"aa", age:50})           |
| $and   | 显式逻辑与 | {$and: [{<key1>:<value1>, <key2>:<value2>, ...}]} | db.xx.find({$and: [{name:"aa", age:50}]}) |
| $or    | 逻辑或     | {$or: [{<key1>:<value1>, <key2>:<value2>, ...}]}  | db.xx.find({$or: [{name:"aa", age:50}]})  |
| $not   | 取反       | {<key>: {$not:{condition1, condition2, ...}}}     | db.xx.find(age:{$not:{$lte:30,$gte:20}})  |

```

AND 和 OR 一起使用
找出姓名为aa或年龄小于50的男性:
db.xx.find({sex:"男", $or: [{name:"aa", age:{$lt:50}}]})
```



#### 数组相关表达式

| 表达式     | 语义                             | 语法                          | 示例                                           |
| :--------- | :------------------------------- | ----------------------------- | ---------------------------------------------- |
| $in        | 验证key的值在value_array里的项   | {<key>:{$in: <value_array>}}  | db.student.find({age:{$in: [50,20,6]}});       |
| $nin       | 验证key的值不在value_array里的项 | {<key>:{$nin: <value_array>}} | db.student.find({age:{$nin: [50,20,6]}});      |
| $all       | 验证value_array都在key里的项     | {<key>:{$all: <value_array>}} | db.student.find({ailas:{$all: ["A3", "A2"]}}); |
| $size      | 验证数组key大小为value的项       | {<key>: {$size:<value>}}      | db.student.find({ailas:{$size: 2}});           |
| $elemMatch | 匹配数组内的元素                 |                               |                                                |

官网上说不能用来匹配一个范围内的元素，如果想找$size<5之类的，他们建议创建一个字段来保存元素的数量。

#### 字段相关表达式 

| 操作符  | 语义                  | 语法                      | 示例                                                         |
| :------ | :-------------------- | ------------------------- | ------------------------------------------------------------ |
| $exists | 验证字段key存在       | {<key>: {$exists: true}}  | db.student.find({score:{$exists:true}});                     |
| $type   | 验证字段key的数据类型 | {<key>: {$type: type_id}} | db.things.find({a: {$type: 16}}); 	// matches if a is an int |

#### $mod取模运算

```
假设age模12余0表示本命年:
db.things.find({{age:{$mod:[12,0]}});
```



 #### 模糊查询与正则表达式

查询 title 包含"教"字的文档：

```
查询 title 包含"教"字的文档：
db.col.find({title:/教/})

查询 title 字段以"教"字开头的文档：
db.col.find({title:/^教/})

查询 title字段以"教"字结尾的文档：
db.col.find({title:/教$/})

查询 title字段不以"教"字结尾的文档：
db.col.find({title:{$not:/教$/})

```

#### $where修饰符

它是一个非常强大的修饰符，但强大的背后也意味着有风险存在。它可以让我们在条件里使用javascript的方法来进行复杂查询。我们先来看一个最简单的例子，现在要查询年龄大于30岁的人员。

```
db.workmate.find(
    {$where:"this.age>30"},
    {name:true,age:true,_id:false}
)
```

这里的this指向的是workmate（查询集合）本身。这样我们就可以在程序中随意调用。虽然强大和灵活，但是这种查询对于数据库的压力和安全性都会变重，所以在**工作中尽量减少$where修饰符的使用**。

### 结果操作

- query：这个就是查询条件，MongoDB默认的第一个参数。
- fields：（返回内容）查询出来后显示的结果样式，可以用true和false控制是否显示。
- limit：返回的数量，后边跟数字，控制每次查询返回的结果数量。0表示所有
- skip:跳过多少个显示，和limit结合可以实现分页。
- sort：排序方式，从小到大排序使用1，从大到小排序使用-1。
- count():结果集总数。
- pretty():结果格式化输出。

```

分页Demo：
db.workmate.find({},{name:true,age:true,_id:false}).limit(0).skip(2).sort({age:1});

```



### 常用语句备注

```
数组定位:
比如我们现在要修改xiaoWang的第三个兴趣(interest)为编码（Code），注意这里的计数是从0开始的。
db.workmate.update({name:'xiaoWang'},{$set:{"interest.2":"Code"}})

喜欢看电影而且看书的人员
db.workmate.find({interest:{$all:["看电影","看书"]}},{name:1,interest:1,age:1,_id:0})

喜欢看电影或看书的人员, 类似$or
db.workmate.find({interest:{$in:["看电影","看书"]}},{name:1,interest:1,age:1,_id:0})

喜欢看电影的人员
db.workmate.find({interest:'看电影'},{name:1,interest:1,age:1,_id:0} )
注意'看电影'不要用中括号括起来, 因为加上中括号就相当于完全匹配了

比如现在我们知道了一个人的爱好是’画画’,’聚会’,’看电影’，但我们不知道是谁，这时候我们就可以使用最简单的数组查询
db.workmate.find({interest:['画画','聚会','看电影']}, {name:1,interest:1,age:1,_id:0})

```



#### 

