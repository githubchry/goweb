package logics

import (
	"context"
	"log"
)

type EventUploadServiceImpl struct{}

func (p *EventUploadServiceImpl) EventUpload(ctx context.Context, args *EventReq, ) (*EventRsp, error) {
	// 0.grpc收到报警事件,
	// 1.根据时间结构体imgurl字段, 使用http get取图片数据
	// 2.图片数据转为base64
	// 3.post base64到算法模块 得到特征数据
	// 4.特征数据发送到milus, 得到id
	// 5.id+特征+报警数据保存到mongo


	log.Println(args);
	rsp := &EventRsp{Message: "sucess"}
	return rsp, nil
}



/*
1. 通过web api post图片到服务器
2. 服务器协程A接受 生成事件 推送到kafka 可以结构体丢topicA  图片丢topicB
3. 服务器协程B从kafka topicAB取事件, 积攒16个或超时后, 调用C++库处理(先不搞C++ 模拟出来即可)
4. 服务器从C++库得到处理结果,可能是异步回调方式,也可能是同步方式, 把结果推送到kafka topicC
5. 服务器协程C从kafka topicC获取到结果, 一个个地返回给谁??
*/