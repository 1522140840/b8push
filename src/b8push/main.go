package main

import (
	al "b8push/log"
	"fmt"
	"os"
	"runtime"
	"os/signal"
)

var log = al.InitLog("default")
func main()  {
	defer recoverFromError()
	//设置最大能够同时执行的协程数,在环境变量里配置
	// runtime.GOMAXPROCS(8)
	log.Infof("runtime.GOMAXPROCS %d", runtime.GOMAXPROCS(0))

	//注册系统中断，关闭信息
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)


	//启动websocket服务
	go StartServer()


	<-c
	log.Infof("Push server stopped\n")
	log.Flush()
}

func recoverFromError() {
	if err := recover(); err != nil {
		fmt.Fprintln(os.Stderr, "Recover error:", err)
	}
}
