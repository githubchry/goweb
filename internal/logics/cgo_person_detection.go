package logics

//#cgo CFLAGS: -I${SRCDIR}/../../cgo
//#cgo LDFLAGS: -L${SRCDIR}/../../cgo/person_detection -lperson_detection -L${SRCDIR}/../../cgo/3rdparty/lib -lprotobuf -lstdc++
// //编译时GCC会自动找到libnumber.a或libnumber.so进行链接
//#include <person_detection/person_detection.h>
import "C"
import (
	"github.com/githubchry/goweb/internal/logics/protos"
	"github.com/golang/protobuf/proto"
	"log"
	"unsafe"
)

func PersonDetection(input *protos.PersonDetectionInput) *protos.PersonDetectionOutput {

	//把logics.UserLoginRsp结构体转成protobuf二进制数据
	indata, _ := proto.Marshal(input)

	outdata := make([]byte, 1024)

	log.Println("input protobuf:", len(indata))

	// 这里不拷贝img了 直接把img的地址共享给C语言函数
	ret := int(C.person_detection((*C.char)(unsafe.Pointer(&indata[0])), C.int(len(indata)), (*C.char)(unsafe.Pointer(&outdata[0])), C.int(len(outdata))))

	log.Println("return protobuf:", ret)

	var output protos.PersonDetectionOutput
	if err := proto.Unmarshal(outdata, &output); err != nil {
		log.Println("Failed to parse protobuf:", err)
		return nil
	}

	return &output
}