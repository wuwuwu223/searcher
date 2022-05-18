package model

import "gorm.io/gorm"

// Data 数据存储结构体
type Data struct {
	gorm.Model
	Url     string `json:"url"`
	Caption string `json:"caption"`
	Kws     []Kw   `gorm:"foreignKey:DataId" json:"kws"`
}

// DataResult 搜索结果结构体
type DataResult struct {
	Url     string `json:"url"`
	Caption string `json:"caption"`
}
