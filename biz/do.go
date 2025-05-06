package biz

/*
#include <stdio.h>
#include <stdlib.h>

static void myprint(char* s) {
	printf("%s (printed in C)\n", s);
}
*/
import "C"
import "unsafe"

func DoC() {
	cs := C.CString("A string from Go side")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
