/**
 * pbf: proto文件路径, 如"proto/user.proto"
 * reqtype: proto文件里面的请求结构体, 如"logics.UserLoginReq"
 * reqjson: reqtype对应的json结构体数据
 * rsptype: proto文件里面的应答结构体, 如"logics.UserLoginRsq"
 * handle: 处理返回函数
 * */
function pbcall(pbf, req_type, rsp_type, url, req_json, handle) {

    // protobuf.load会把proto文件解析成一个树状结构 得到root节点
    protobuf.load(pbf, function (err, root) {
        if (err) throw err;

        console.log(req_json)
        // 获取消息类型Obtain a message type
        var pbreq = root.lookupType(req_type);
        // 验证有效负载(如可能不完整或无效 比如password, proto定义是string类型, 如果payload设为int的数值,就会返回错误)
        if (pbreq.verify(req_json)) throw Error(errMsg);

        // 根据payload创建消息体
        var req_type_obj = pbreq.create(req_json);
        console.log(req_type_obj);

        // Encode a message to an Uint8Array (browser) or Buffer (node)
        var req_type_obj_buf = pbreq.encode(req_type_obj).finish();
        const options = { body: req_type_obj_buf, method: "POST" };

        // 发出http post请求, 等待响应
        fetch(url, options).then(function (response) {
            response.arrayBuffer().then(function (rsp_type_obj_buf) {
                // 获取消息类型Obtain a message type
                var pbrsp = root.lookupType(rsp_type);
                // 获取消息类型Obtain a message type
                var rsp_type_obj = pbrsp.decode(new Uint8Array(rsp_type_obj_buf))
                console.log(rsp_type_obj);
                handle(rsp_type_obj);
            })
        });
    });
}