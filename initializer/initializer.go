package initializer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"os"
	"searcher/global"
	"searcher/search/model"
	"searcher/search/utils"
	"searcher/web"
	"strings"
	"sync"
	"time"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigType("json")
	configFileName := fmt.Sprintf("config.json")
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.Config); err != nil {
		panic(err)
	}

}
func InitDb() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             3 * time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Silent,   // 日志级别
			IgnoreRecordNotFoundError: false,           // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,           // 禁用彩色打印
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
		global.Config.Mysql.Username,
		global.Config.Mysql.Password,
		global.Config.Mysql.Address,
		global.Config.Mysql.Database)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true, //开启缓存
		SkipDefaultTransaction: true, //关闭默认事务
		Logger:                 newLogger,
	})
	sqlDB, err := Db.DB()
	sqlDB.SetMaxIdleConns(100)  //空闲连接数
	sqlDB.SetMaxOpenConns(1000) //最大连接数
	sqlDB.SetConnMaxLifetime(time.Minute)
	if err != nil {
		fmt.Println(err)
	}
	Db.AutoMigrate(&model.Data{}, &model.Kw{}) //创建数据表
	global.Db = Db
}

func InitSegmenter() {
	global.Seg.LoadDictionary("dictionary.txt") //加载字典
}

func InitHttpServer(r *gin.Engine) {
	apiRouter := r.Group("/api/v1")
	apiRouter.GET("search", web.Search)
}

func InitData() {
	fileInfoList, err := ioutil.ReadDir("resources")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for i := range fileInfoList {
		if strings.Contains(fileInfoList[i].Name(), "csv") {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				log.Println("正在导入", fileInfoList[i].Name())
				utils.ImportCsv("resources/" + fileInfoList[i].Name())
				log.Println("导入完成", fileInfoList[i].Name())
			}(i)
		}
		if i != 0 && i%2 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}
