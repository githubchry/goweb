package algorithm

/*
#include <stdio.h>
int add(int a, int b) {
    return a+b;
}
*/
//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
// //编译时GCC会自动找到libnumber.a或libnumber.so进行链接
//#include "number.h"
import "C"


// 使用cgo模块之前得保证PATH环境里面有gcc  比如win10下添加PATH路径=> D:\Qt5.14.2\Tools\mingw730_64\bin

func TestAdd(a int, b int) int {
	return int(C.add(C.int(a), C.int(b)))
}

func TestAddMod(a int, b int, c int) int {
	return int(C.number_add_mod(C.int(a), C.int(b), C.int(c)))
}