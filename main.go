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
	"strings"
	"sync"
	"time"
)

type Task struct {
}

func (t Task) SimulateWork(id int) {
	fmt.Println("work")
}

var maxGoroutine = 100

var allocOpen bool

func init() {
	mg, _ := strconv.Atoi(os.Getenv("MAX_GOROUTINE"))
	if mg > 0 {
		maxGoroutine = mg
	}
	log.Printf("Max Goroutine: %d", maxGoroutine)
	if t := strings.ToLower(os.Getenv("ALLOC_OPEN")); t == "1" || t == "t" || t == "true" {
		allocOpen = true
	}
	log.Printf("Alloc Open: %t", allocOpen)
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
	var (
		counter      int
		counterMutex sync.Mutex
	)
	http.HandleFunc("/mutex", func(w http.ResponseWriter, r *http.Request) {
		counterMutex.Lock()
		defer counterMutex.Unlock()

		counter++
		fmt.Fprintf(w, "Counter: %d\n", counter)
	})
	server := &http.Server{
		Addr: ":9090",
	}
	http.HandleFunc("/busy/", func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(r.URL.Path, "/")
		if len(paths) < 3 {
			http.Error(w, "xxx", http.StatusBadRequest)
		}
		atoi, _ := strconv.Atoi(paths[2])
		if atoi == 0 {
			atoi = 10
		}
		var wg sync.WaitGroup
		for i := 0; i < atoi; i++ {
			wg.Add(1)
			curr <- struct{}{}
			go func() {
				defer wg.Done()
				doBusy()
				<-curr
			}()
		}
		wg.Wait()
	})
	if allocOpen {
		go simulateAlloc()
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type largeObject struct {
	data [1024 * 1024]byte
}

func simulateAlloc() {
	for {
		largeObject := &largeObject{}

		fmt.Printf("Allocated large object with address: %p\n", largeObject)

		time.Sleep(time.Second * 3)
	}
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

func doBusy() {
	for i := 0; i < 128; i++ {
		doThing(i)
	}
}
