package global

import (
	"github.com/huichen/sego"
	"gorm.io/gorm"
	"searcher/config"
)

var Db *gorm.DB
var Seg sego.Segmenter
var Config *config.Config
