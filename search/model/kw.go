package model

import "gorm.io/gorm"

type Kw struct {
	gorm.Model
	Word   string `json:"word" gorm:"index:idx_word_dataid,unique"`
	DataId uint   `json:"data_id" gorm:"index:idx_word_dataid,unique"`
}
