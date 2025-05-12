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
func uppercaseByC(ctx context.Context, str string, args ...interface{}) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	// 默认睡眠时间为0，如果传入可变参数并类型断言为int成功，则获取此值作为睡眠时间
	sleepTime := 0
	if len(args) > 0 {
		if st, ok := args[0].(int); ok {
			sleepTime = st
		}
	}

	C.to_uppercase(cStr, C.int(sleepTime))

	return C.GoString(cStr), nil
}
