package biz

import (
	"context"
	"errors"
)

func DoMalloc(ctx context.Context, args interface{}) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	conf, ok := args.(MallocConf)
	if !ok {
		return errors.New("?")
	}
	var f func(int, int) interface{}
	switch conf.AllocType {
	case 1:
		f = doPreallocate
	case 2:
		f = doPreallocate2
	default:
		f = withoutPreallocate
	}
	_ = f(conf.Capacity, conf.Bytes)
	return nil
}

type MallocConf struct {
	AllocType int
	Capacity  int
	Bytes     int
}

func doPreallocate2(capa, bytes int) interface{} {
	data := make([][]byte, capa)
	for j := 0; j < capa; j++ {
		newBytes := make([]byte, bytes)
		data[j] = newBytes
	}
	return data
}

func doPreallocate(capa, bytes int) interface{} {
	data := make([][]byte, 0, capa)
	for j := 0; j < capa; j++ {
		newBytes := make([]byte, bytes)
		data = append(data, newBytes)
	}
	return data
}

func withoutPreallocate(capa, bytes int) interface{} {
	var data [][]byte
	for j := 0; j < capa; j++ {
		newBytes := make([]byte, bytes)
		data = append(data, newBytes)
	}
	return data
}
