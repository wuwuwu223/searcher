package initializer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"searcher/global"
	"searcher/search/model"
	"searcher/web"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
		global.Config.Mysql.Username,
		global.Config.Mysql.Password,
		global.Config.Mysql.Address,
		global.Config.Mysql.Database)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true, //开启缓存
		SkipDefaultTransaction: true, //关闭默认事务
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
