package search

import (
	"gorm.io/gorm/clause"
	"searcher/global"
	"searcher/search/model"
	"searcher/search/utils"
	"sort"
	"strings"
)

func Search(row string) []model.DataResult {
	arr := utils.SplitStr(row) //分词

	var result []model.Kw
	for i := range arr { //遍历查询分词结果
		var alist []model.Kw
		//global.Db.Where("word=?", arr[i]).FindInBatches(&alist, 500, func(tx *gorm.DB, batch int) error {
		//	result = append(result, alist...)
		//	return nil
		//})
		global.Db.Where("word=?", arr[i]).Find(&alist)
		result = append(result, alist...)
	}
	mapping := make(map[uint]int, len(result)) //分析词频
	for j := range result {
		token := strings.Split(result[j].Word, "/")
		cx := token[1]
		if cx == "nr" {
			mapping[result[j].DataId]++
		}
		mapping[result[j].DataId]++
	}
	var sdatas sDatas
	for k, v := range mapping {
		var sdata sData
		sdata.Id = k
		sdata.Count = v
		sdatas = append(sdatas, sdata)
	}
	sort.Sort(&sdatas) //词频排序
	var values []uint
	for i := range sdatas {
		values = append(values, sdatas[i].Id)
	}
	var list []model.DataResult
	global.Db.Model(&model.Data{}).Select("url,caption").Where("id in (?)", values).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{values}, WithoutParentheses: true},
	}).Find(&list) //查询数据
	return list //返回数据
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
