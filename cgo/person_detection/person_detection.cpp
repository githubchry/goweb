#include "person_detection.h"
#include "person_detection_wrapper.h"

int person_detection(char *input, int inlen, char *output, int outlen)
{
	return person_detection_wrapper(input, inlen, output, outlen);
}

/*
静态链接库生成
gcc -c person_detection.cpp person_detection_wrapper.cpp algorithm.pb.cc -I../3rdparty/include
ar rcs libperson_detection.a person_detection.o person_detection_wrapper.o algorithm.pb.o



linux动态链接库生成
gcc -shared -o libnumber.so number.c

windows动态库生成 库文件不强制要求lib开头
gcc number.c -shared -o number.dll -Wl,--out-implib,number.a
number.dll文件需放在运行目录

额外生成.def文件
gcc number.c -shared -o number.dll -Wl,--output-def,number.def,--out-implib,number.a


//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
// //编译时GCC会自动找到libnumber.a或libnumber.so进行链接  windows下不强制要求lib开头
*/