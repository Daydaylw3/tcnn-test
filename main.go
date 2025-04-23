package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
)

type Task struct {
}

func (t Task) SimulateWork(id int) {
	fmt.Println("work")
}

var maxGoroutine = 100

func init() {
	mg, _ := strconv.Atoi(os.Getenv("MAX_GOROUTINE"))
	if mg > 0 {
		maxGoroutine = mg
	}
}

func main() {
	// 开启对阻塞和锁调用的跟踪
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	
	curr := make(chan struct{}, maxGoroutine)
	task := &Task{}
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
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
	})
	server := &http.Server{
		Addr: ":9090",
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
