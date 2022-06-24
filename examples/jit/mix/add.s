#include "textflag.h"

// 静态基地址(static-base) 指针
//    |
//                  |         add函数入参+返回值总大小
//                  |               |
// TEXT pkgname·add(SB),NOSPLIT,$16-24
//      |      |                |
// 函数所属包名  函数名          add函数栈帧大小
//

TEXT ·Add(SB), NOSPLIT, $0-24
    ADDQ BX, AX
    RET
