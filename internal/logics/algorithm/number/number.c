#include "number.h"

int number_add_mod(int a, int b, int mod) {
    return (a+b)%mod;
}

/*
gcc -c number.c
得到number.o文件, -c表示编译为目标文件

ar rcs libnumber.a number.o
得到libnumber.a文件
*/