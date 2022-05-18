package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"searcher/global"
	"searcher/initializer"
)

func main() {
	needInit := false
	flag.BoolVar(&needInit, "init", false, "是否初始化数据库")
	flag.Parse()
	initializer.InitConfig()
	initializer.InitDb()
	initializer.InitSegmenter()
	if needInit { //导入resources下的csv文件
		initializer.InitData()
		return //导入后关闭程序
	}
	r := gin.Default()

	initializer.InitHttpServer(r)
	err := r.Run(fmt.Sprintf("%s:%d", global.Config.HttpServer.ListenIp, global.Config.HttpServer.ListenPort))
	if err != nil {
		panic(err)
	}
}
