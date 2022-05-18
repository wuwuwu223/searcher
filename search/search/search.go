package search

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"searcher/global"
	"searcher/search/model"
	"searcher/search/utils"
	"sort"
)

func Search(row string) []model.DataResult {
	arr := utils.SplitStr(row)
	mapping := make(map[uint]int)
	var result []model.Kw
	for i := range arr {
		var alist []model.Kw
		global.Db.Where("word=?", arr[i]).FindInBatches(&alist, 10000, func(tx *gorm.DB, batch int) error {
			result = append(result, alist...)
			return nil
		})
		for j := range result {
			mapping[result[j].DataId]++
		}
	}
	var sdatas sDatas
	for k, v := range mapping {
		var sdata sData
		sdata.Id = k
		sdata.Count = v
		sdatas = append(sdatas, sdata)
	}
	sort.Sort(&sdatas)
	var values []uint
	for i := range sdatas {
		values = append(values, sdatas[i].Id)
	}
	var list []model.DataResult
	global.Db.Model(&model.Data{}).Select("url,caption").Where("id in (?)", values).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{values}, WithoutParentheses: true},
	}).Find(&list)
	return list
}

type sData struct {
	Id    uint
	Count int
}

type sDatas []sData

func (s sDatas) Less(i, j int) bool {
	return s[i].Count > s[j].Count
}

func (s sDatas) Swap(i, j int) {
	//TODO implement me
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func (s sDatas) Len() int { return len(s) }
