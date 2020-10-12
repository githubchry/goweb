package algorithm

/*
#include <stdio.h>
int add(int a, int b) {
    return a+b;
}
*/
import "C"


// 使用cgo模块之前得保证PATH环境里面有gcc  比如win10下添加PATH路径=> D:\Qt5.14.2\Tools\mingw730_64\bin

func TestAdd(a int, b int) int {
	return int(C.add(C.int(a), C.int(b)))
}