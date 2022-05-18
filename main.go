package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"searcher/global"
	"searcher/initializer"
	"searcher/search/utils"
)

func main() {
	needInit := false
	flag.BoolVar(&needInit, "init", false, "是否初始化数据库")
	flag.Parse()
	initializer.InitConfig()
	initializer.InitDb()
	initializer.InitSegmenter()
	if needInit { //如果需要初始化就导入这个文件
		utils.ImportCsv("wukong50k_release.csv")
	}
	r := gin.Default()

	initializer.InitHttpServer(r)
	err := r.Run(fmt.Sprintf("%s:%d", global.Config.HttpServer.ListenIp, global.Config.HttpServer.ListenPort))
	if err != nil {
		panic(err)
	}
}
