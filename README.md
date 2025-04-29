# tcnn-test

## 编译

```shell
## 生成 tcnn-test
go build
## 指定文件名 xxxx
go build -o xxxx
```

## 启动

```shell
## 默认
## 最高同时 goroutine = 100
## 不开启 每3s分配 1M 内存(尽可能触发gc)
## 端口 :9090
./tcnn-test

## 进行配置
## 最高同时 goroutine = 10
## 开启 每3s分配 1M 内存(尽可能触发gc)
## 端口 :8080(*不能省略前面的冒号)
MAX_GOROUTINE=10 ALLOC_OPEN=t ENDPOINT=:8080 ./tcnn-test

## 如果要进行 numa 邦核
MAX_GOROUTINE=10 ALLOC_OPEN=t ENDPOINT=:8080 numactl -C 24-31 ./tcnn-test
```

## 跑压测

+ 当前使用 `endpoint/busy/1000` 来进行测试, 其中 `1000` 为任务总数

## 获取 pprof 数据

+ 获取过去`10s`(可配置)的 trace: `curl -o trace.trace http://{endpoint}/debug/pprof/trace?secondes=10`
+ 获取过去`10s`(可配置)的 pprof: `curl -o pprof.pprof http://{endpoint}/debug/pprof/profile?secondes=10`

## 启动 web 服务进行查看数据

+ trace: `go tool trace -http=:8080 trace.trace`, 启动在 `localhost:8080`
+ pprof: `go tool pprof -http=:8080 pprof.pprof`, 启动在 `localhost:8080`

