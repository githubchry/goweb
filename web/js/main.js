let sendChannel = null;
let remoteVideo = document.getElementById('remoteVideo');
let pc = new RTCPeerConnection(null);

// onnegotiationneeded在要求sesssion协商时发生
pc.onnegotiationneeded = async function() {
    // 创建并设置本地SDP(会话描述协议Session Description Protocol) 
    let offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    // 获取远端SDP
    getRemoteSdp();
}

// onaddstream在远程数据流到达时发生，将数据流装载到video中
pc.onaddstream = function(event) {
    remoteVideo.srcObject = event.stream;
}

async function getRemoteSdp() {
    // const url = 'http://172.17.58.119:8083/recive';
    const url = '/api/swapsdp';
    const options  = {
        headers: { "content-type": "application/json; charset=UTF-8" },
        // 将请求javascript object转换为JSON字符串
        body: btoa(pc.localDescription.sdp),
        method: "POST"
    }
    
    console.log(btoa(pc.localDescription.sdp))
    // 发出http POST 请求
    let response = await fetch(url, options)
    let remoteSdp = await response.text()
    console.log(remoteSdp)

    pc.setRemoteDescription(new RTCSessionDescription({
        type: 'answer',
        sdp: atob(remoteSdp)
    }))
}
  
async function getCodecInfo() {
    // const url = 'http://172.17.58.119:8083/codec/demo1';
    const url = '/api/getcodec';
    const options  = {
        headers: { "content-type": "application/json; charset=UTF-8" },
        body: "rtsp://172.24.69.55:8554/slamtv60.264",
        method: "POST"
    }

    // 发出http post 请求
    let response = await fetch(url, options)

    // 等待响应，并将响应作为JSON字符串反序列化为javascript object
    let js = await response.json()
    console.log(js)

    // 为pc添加音视频收发器, 即可发送数据, 亦可接收数据; 将触发onnegotiationneeded
    pc.addTransceiver('video', {'direction': 'sendrecv'});
    if (js.length > 1) {
        pc.addTransceiver('audio', {'direction': 'sendrecv'});
    }

    //send ping becouse PION not handle RTCSessionDescription.close()
    sendChannel = pc.createDataChannel('foo');

    sendChannel.onclose = () => console.log('sendChannel has closed');
    
    sendChannel.onopen = () => {
        console.log('sendChannel has opened');
        sendChannel.send('ping');
        setInterval(() => {
            sendChannel.send('ping');
        }, 1000)
    }

    sendChannel.onmessage = e => console.log(`Message from DataChannel '${sendChannel.label}' payload '${e.data}'`);
    
}
  
window.onload = function() {
    getCodecInfo();
}