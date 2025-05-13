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
	"context"
	"unsafe"
)

// 封装了调用 C 来转换大写的方法
// 可变参数用于传入额外参数，如睡眠时间
func uppercaseByC(ctx context.Context, str string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	C.to_uppercase(cStr)

	return C.GoString(cStr), nil
}
