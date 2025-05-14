package biz

import (
	"context"
	"strconv"
)

func DoMalloc(ctx context.Context, args ...string) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	a := getArg(args)
	multi := 10000
	if len(args) > 3 {
		i, err := strconv.Atoi(args[2])
		if err == nil && i > 0 {
			multi = i
		}
	}
	switch a.preAlloc {
	case 0:
		for i := 0; i < multi; i++ {
			_ = withoutPreallocate(a.capacity)
		}
	case 1:
		for i := 0; i < multi; i++ {
			_ = doPreallocate(a.capacity)
		}
	case 2:
		for i := 0; i < multi; i++ {
			_ = doPreallocate2(a.capacity)
		}
	}
	return nil
}

type arg struct {
	capacity int
	preAlloc int
}

func getArg(args []string) arg {
	defCap := 10000
	defPre := 0
	switch len(args) {
	case 0:
		return arg{capacity: defCap}
	case 1:
		i, err := strconv.Atoi(args[0])
		if err != nil || i <= 0 {
			return arg{capacity: defCap}
		}
		return arg{capacity: i}
	default:
		i, err := strconv.Atoi(args[0])
		if err == nil || i > 0 {
			defCap = i
		}
		p, err := strconv.Atoi(args[1])
		if err == nil || i > 0 {
			defPre = p
		}
		return arg{capacity: defCap, preAlloc: defPre}
	}
}

func doPreallocate2(capa int) interface{} {
	data := make([][]byte, capa)
	for j := 0; j < capa; j++ {
		newBytes := make([]byte, 256)
		data[j] = newBytes
	}
	return data
}

func doPreallocate(capa int) interface{} {
	data := make([][]byte, 0, capa)
	for j := 0; j < capa; j++ {
		newBytes := make([]byte, 256)
		data = append(data, newBytes)
	}
	return data
}

func withoutPreallocate(capa int) interface{} {
	var data [][]byte
	for j := 0; j < capa; j++ {
		newBytes := make([]byte, 256)
		data = append(data, newBytes)
	}
	return data
}
