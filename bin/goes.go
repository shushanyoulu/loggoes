package main

import (
	"runtime"

	"os"

	"github.com/Shopify/sarama"

	_ "net/http/pprof"
)

var offLine = make(map[string]onAndOff) //退出信息
var (
	zookeeper      = zookeeperAddr()
	zookeeperNodes []string
	m, k           []string
)

// var loginInfoIdx = esLoginInfoIndex() //登陆详细信息
// var logOutIdx = esLogoutIndex()

// var parallelNum = tcount()
var goesl = goeslog() //配置日志字段

var logi, logj uint32
var put2es = putToEsBuff()

func init() {
	sarama.Logger = goeslog()

}

var listenGetNum uint32          //用以记录系统已经收到了多少条数据
var over = make(chan string, 10) // 运行结束通知
var kafkaConfigInfo = ktopic()   //读取kafka消费者信息
func main() {
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
	// go func() {
	// 	http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
	// 		num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
	// 		w.Write([]byte(num))
	// 	})
	// 	http.ListenAndServe("localhost:6060", nil)
	// 	glog.Info("goroutine stats and pprof listen on 6060")
	// }()
	runtime.GOMAXPROCS(uCPUNum()) //配置程序可用cpu数量
	go tikers()
	go judgeZeroTime()
	analysisTopicGroup(kafkaConfigInfo)
	goesl.Println(<-over)
	os.Exit(0)
} // main
