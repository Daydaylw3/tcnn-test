//go:generate gcc -c -o libsrc/libtool.a libsrc/tool.c
package biz

/*
#cgo CFLAGS: -g -O2 -I${SRCDIR}/libsrc
#cgo LDFLAGS: -L${SRCDIR}/libsrc -ltool

#include <stdlib.h>
#include "tool.h"
*/
import "C"
import (
	"unsafe"
)

// 封装了调用 C 来转换大写的方法
func uppercaseByC(str string) string {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	C.to_uppercase(cStr)

	return C.GoString(cStr)
}
