#include "number.h"

int number_add_mod(int a, int b, int mod) {
    return (a + b) % mod;
}

/*
静态链接库生成
gcc -c number.c
得到number.o文件, -c表示编译为目标文件

ar rcs libnumber.a number.o
得到libnumber.a文件


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