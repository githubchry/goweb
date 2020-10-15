package logics

/*
#include <stdio.h>
#include <fcntl.h>      // O_CREAT | O_RDWR | O_BINARY
#include <stdlib.h>     // defer C.free
int add(int a, int b) {
    return a+b;
}

int saveimg(char *filename, char *ptr, int size) {
    int fd = open(filename, O_CREAT | O_RDWR, "0777");
	if (fd < 0) return -1;
	int ret = write(fd, ptr, size);
	if (ret < size) return -2;
	close(fd);
    return ret;
}

int save_multipath_img(char **ptr, int *size, int num) {


    return 0;
}
*/
//#cgo CFLAGS: -I./algorithm/number
//#cgo LDFLAGS: -L${SRCDIR}/algorithm/number -lnumber
// //编译时GCC会自动找到libnumber.a或libnumber.so进行链接
//#include "algorithm/number/number.h"
import "C"
import (
	"unsafe"
)

// 使用cgo模块之前得保证PATH环境里面有gcc  比如win10下添加PATH路径=> D:\Qt5.14.2\Tools\mingw730_64\bin

func TestAdd(a int, b int) int {
	return int(C.add(C.int(a), C.int(b)))
}

func TestAddMod(a, b, c int) int {
	return int(C.number_add_mod(C.int(a), C.int(b), C.int(c)))
}

func TestSaveImg(filename string, img []byte) int {
	// go的string是一个带len的结构体 c语言的char*是以\0结尾的数组
	path := C.CString(filename)		 // 这里根据go string在堆上新建适用于C语言的char* 需要手动释放
	defer C.free(unsafe.Pointer(path))

	//var CBytesPtr unsafe.Pointer = C.CBytes(img)	// 这里根据go []byte在堆上新建适用于C语言的char* 需要手动释放
	//defer C.free(CBytesPtr)

	// 这里不拷贝img了 直接把img的地址共享给C语言函数
	return int(C.saveimg(path, (*C.char)(unsafe.Pointer(&img[0])), C.int(len(img))))
}

func TestSaveImgs(imgs [][]byte) int {
	//var CBytesPtr unsafe.Pointer = C.CBytes(img)	// 这里根据go []byte在堆上新建适用于C语言的char* 需要手动释放
	//defer C.free(CBytesPtr)

	var imgnum int
	var imgssize []int

	imgnum = len(imgs)
	for i := 0; i < imgnum; i++ {
		imgssize[i] = len(imgs[i])
	}

	// 这里不拷贝img了 直接把img的地址共享给C语言函数
	return int(C.save_multipath_img((**C.char)(unsafe.Pointer(&imgs[0][0])), (*C.int)(unsafe.Pointer(&imgssize[0])),C.int(imgnum)))
}