package model

import "gorm.io/gorm"

type Kw struct {
	gorm.Model
	Word   string `json:"word"`
	DataId uint   `json:"data_id"`
}
