package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"tcnn-test/biz"
)

var curr chan struct{}

var task = &Task{}

func work(_ http.ResponseWriter, _ *http.Request) {
	var wg sync.WaitGroup
	curr <- struct{}{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		method := reflect.ValueOf(task).MethodByName("SimulateWork")
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(1)})
		}
		<-curr
	}()
	wg.Wait()
}

func mutexFunc() http.HandlerFunc {
	var (
		counter      int
		counterMutex sync.Mutex
	)
	return func(w http.ResponseWriter, r *http.Request) {
		counterMutex.Lock()
		defer counterMutex.Unlock()

		counter++
		fmt.Fprintf(w, "Counter: %d\n", counter)
	}
}

func busy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskC := getTaskCount(r, 10000)
		start := time.Now()

		if done, err := assignTask(r.Context(), taskC, doBusy); err != nil {
			cost := time.Since(start).Round(time.Millisecond)
			log.Printf("busy  job canceled, %10d tasks done, cost: %s", done, cost)
			return
		}
		cost := time.Since(start)
		log.Printf("busy  job finish, %10d tasks done, cost: %s", taskC, cost.Round(time.Second))
		fmt.Fprintf(w, "busy  job finish, %10d tasks done, cost: %s\n", taskC, cost.Round(time.Millisecond))
	}
}

func getTaskCount(r *http.Request, def int) int {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) >= 3 {
		if i, _ := strconv.Atoi(paths[2]); i > 0 {
			return i
		}
	}
	return def
}

func assignTask(ctx context.Context, taskC int, do func(context.Context, ...string) error, args ...string) (int, error) {
	var wg sync.WaitGroup
	for i := 0; i < taskC; i++ {
		wg.Add(1)
		curr <- struct{}{}
		select {
		case <-ctx.Done():
			return i + 1, ctx.Err()
		default:
			go func() {
				defer wg.Done()
				do(ctx, args...)
				<-curr
			}()
		}
	}
	wg.Wait()
	return taskC, ctx.Err()
}

func doBusy(ctx context.Context, _ ...string) error {
	for i := 0; i < 128; i++ {
		select {
		case <-ctx.Done():
			// 如果上下文被取消，退出并返回错误
			return ctx.Err()
		default:
			// 没有取消信号，进行正常的工作
			doThing(i)
		}
	}
	return nil
}

var mutex sync.Mutex
var funcTypes []reflect.Type

func doThing(n int) reflect.Type {
	mutex.Lock()
	defer mutex.Unlock()
	if n >= len(funcTypes) {
		newFuncTypes := make([]reflect.Type, n+1)
		copy(newFuncTypes, funcTypes)
		funcTypes = newFuncTypes
	}
	if funcTypes[n] != nil {
		return funcTypes[n]
	}
	funcTypes[n] = reflect.StructOf([]reflect.StructField{
		{
			Name: "FuncType",
			Type: reflect.TypeOf(""),
		},
		{
			Name: "Args",
			Type: reflect.ArrayOf(n, reflect.TypeOf("")),
		},
	})
	return funcTypes[n]
}

func malloc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskC := getTaskCount(r, 1000)
		capa := r.URL.Query().Get("cap")
		pre := r.URL.Query().Get("pre")
		multi := r.URL.Query().Get("multi")

		start := time.Now()
		if done, err := assignTask(r.Context(), taskC, biz.DoMalloc, capa, pre, multi); err != nil {
			cost := time.Since(start).Round(time.Millisecond)
			log.Printf("busy  job canceled, %10d tasks done, cost: %s", done, cost)
			return
		}
		cost := time.Since(start)
		log.Printf("busy  job finish, %10d tasks done, cost: %s", taskC, cost.Round(time.Second))
		fmt.Fprintf(w, "busy  job finish, %10d tasks done, cost: %s\n", taskC, cost.Round(time.Millisecond))
	}
}
