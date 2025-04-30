package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Task struct {
}

func (t Task) SimulateWork(id int) {
	fmt.Println("work")
}

var endpoint = ":9090"

var maxGoroutine = 100

var allocOpen bool

func init() {
	if end := os.Getenv("ENDPOINT"); end != "" {
		endpoint = end
	}
	mg, _ := strconv.Atoi(os.Getenv("MAX_GOROUTINE"))
	if mg > 0 {
		maxGoroutine = mg
	}
	log.Printf("Max Goroutine: %d", maxGoroutine)
	curr = make(chan struct{}, maxGoroutine)
	if t := strings.ToLower(os.Getenv("ALLOC_OPEN")); t == "1" || t == "t" || t == "true" {
		allocOpen = true
	}
	log.Printf("Alloc Open: %t", allocOpen)
}

func main() {
	// 开启对阻塞和锁调用的跟踪
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	http.HandleFunc("/work", work)
	http.HandleFunc("/mutex", mutexFunc())
	http.HandleFunc("/busy/", busy())
	http.HandleFunc("/busy", busy())
	server := &http.Server{
		Addr: endpoint,
	}
	if allocOpen {
		go simulateAlloc()
	}
	log.Printf("Server start at: %s", endpoint)
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

		_ = largeObject
		// fmt.Printf("Allocated large object with address: %p\n", largeObject)

		time.Sleep(time.Second * 3)
	}
}
