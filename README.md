# tcnn-test

## 编译

```shell
## 编译 cgo 文件
go generate ./...
## 生成 tcnn-test
go build
## 指定文件名 xxxx
go build -o xxxx
```

## 启动

```shell
## <默认>
## Wiki 文章倍数 = 0
## 不开启数组容量提前分配(该配置只在 busy3 有效果)
## 最高同时 goroutine = 100
## 不开启 每3s分配 1M 内存(尽可能触发gc)
## 端口 :9090
./tcnn-test

## <进行配置>
## Wiki 文章倍数 = 3
## 开启数组容量提前分配(该配置只在 busy3 有效果)
## 最高同时 goroutine = 10
## 开启 每3s分配 1M 内存(尽可能触发gc)
## 端口 :8080(*不能省略前面的冒号)
WIKI_MULTI=3 PRE_MALLOC=t MAX_GOROUTINE=10 ALLOC_OPEN=t ENDPOINT=:8080 ./tcnn-test

## 如果要进行 numa 绑核
MAX_GOROUTINE=10 ALLOC_OPEN=t ENDPOINT=:8080 numactl -C 24-31 ./tcnn-test
```

## 跑压测

+ 当前使用 `<endpoint>/busy/1000` 来进行测试, 其中 `1000` 为任务总数
+ *新增使用 `<endpoint>/busy2/10` 来进行测试, 其中 `10` 为任务总数
+ *新增使用 `<endpoint>/busy3/10` 来进行测试, 其中 `10` 为任务总数

## 获取 pprof 数据

+ 获取过去`10s`(可配置)的 trace: `curl -o trace.trace http://{endpoint}/debug/pprof/trace?secondes=10`
+ 获取过去`10s`(可配置)的 pprof: `curl -o pprof.pprof http://{endpoint}/debug/pprof/profile?secondes=10`

## 启动 web 服务进行查看数据

+ trace: `go tool trace -http=:8080 trace.trace`, 启动在 `localhost:8080`
+ pprof: `go tool pprof -http=:8080 pprof.pprof`, 启动在 `localhost:8080`
