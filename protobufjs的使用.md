

# 编写proto文件
proto/add.proto
```protobuf
syntax = "proto3";      //指明使用proto3语法,否则编译器默认使用proto2的语法
package logics;       //包声明符, 相当于命名空间，产生的类会被包装在C++命名空间中
option go_package = "../internal/logics";
//proto3取消了optional和required限定修饰符，只能使用singular(默认)和repeated
//singular：该字段可以有0个或者1个（但是不能超过1个）。
//repeated：该字段可以重复任意多次（包括0次）。重复的值的顺序会被保留。

//命名规范
//message和enum命名采用驼峰命名方式，大写开头 比如Packet和PacketType
//字段命名采用小写字母加下划线分隔方式 比如packet和packet_type


message AddReq {
  string username = 1;
  string token = 2;
  repeated int32 operand = 3;
}

message AddRsp {
  int32 code = 1;
  string message = 2;
  int64 result = 3;
}
```

**有两种方式在web前端使用protobuf, 具体选择根据个人喜好而定**

# 静态方式
终端运行`protoc --js_out=. add.proto`得到`addreq.js`和`addreq.js`文件

然后在web引用即可

# 动态方式
从 [protobufjs项目地址](https://github.com/protobufjs/protobuf.js) 下载protobuf.js和protobuf.js.map并在web引用:
```javascript
<script src="js/protobuf.js"></script>
```

服务器把proto目录映射出去:
```go
route.PathPrefix("/proto").Handler(http.StripPrefix("/proto", http.FileServer(http.Dir("../proto"))))
```

web代码示例:
``` javascript
// 定义payload数据
const data = {
    username: localStorage.getItem("Username"),
    token:localStorage.getItem("Token"),
    operand: [(Number)(num1), (Number)(num2)]
};

protobuf.load("proto/add.proto", function (err, root) {
    if (err) throw err;

    console.log(data)
    // 获取消息类型Obtain a message type
    var pbreq = root.lookupType("logics.AddReq");
    // 验证有效负载(如可能不完整或无效 比如password, proto定义是string类型, 如果payload设为int的数值,就会返回错误)
    if (pbreq.verify(data)) throw Error(errMsg);

    // 根据payload创建消息体
    var message = pbreq.create(data);
    console.log(message);

    var buffer = pbreq.encode(message).finish();
    const url = '/api/addpost';
    const options = { body: buffer, method: "POST" };

    fetch(url, options).then(function (response) {
        response.arrayBuffer().then(function (buffer) {
            var msg = root.lookupType("logics.AddRsp").decode(new Uint8Array(buffer))
            console.log(msg);
        })
    });
});
```
