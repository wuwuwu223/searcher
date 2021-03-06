package utils

import (
	"encoding/csv"
	"github.com/huichen/sego"
	"io"
	"log"
	"os"
	"searcher/global"
	"searcher/search/model"
	"strings"
	"sync"
)

//ImportCsv 导入csv
func ImportCsv(fileName string) {
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	_, err = r.Read()
	if err != nil {
		log.Fatalf("can not read, err is %+v", err)
	} //跳过首行
	var datas []model.Data
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		var data model.Data
		data.Url = row[0]
		data.Caption = row[1]
		datas = append(datas, data)
	}
	global.Db.CreateInBatches(datas, 1000)
	var wg sync.WaitGroup
	for i := 0; i < len(datas)/100+1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i*100+100 < len(datas) {
				SplitData(datas[i*100 : i*100+100])
			} else {
				SplitData(datas[i*100:])
			}
		}(i)
	}
	wg.Wait()
}

var useless = "dpcueyahokxw"

// SplitData 导入时分词
func SplitData(datas []model.Data) {
	var kws []model.Kw
	for s := range datas {
		text := []byte(datas[s].Caption)
		segments := global.Seg.Segment(text)
		str := sego.SegmentsToString(segments, true)
		arr := strings.Split(str, " ")
		if len(arr) == 0 {
			continue
		}
		existWord := make(map[string]bool)
		for i := range arr {
			if len(arr[i]) > 1 {
				token := strings.Split(arr[i], "/")
				if len(token[1]) < 1 {
					continue
				}
				cx := token[1][0]
				flag := false
				for x := range useless {
					if cx == useless[x] {
						flag = true
					}
				}
				if flag {
					continue
				}
				if existWord[arr[i]] {
					continue
				}
				var kw model.Kw
				kw.DataId = datas[s].ID
				kw.Word = arr[i]
				kws = append(kws, kw)
				existWord[arr[i]] = true
			}
		}
	}
	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	global.Db.Create(kws)
	//log.Println("完成100条", len(kws))
}

// SplitStr 搜索字分词
func SplitStr(row string) []string {
	segments := global.Seg.Segment([]byte(row))
	str := sego.SegmentsToString(segments, true)
	arr := strings.Split(str, " ")
	var result []string
	for i := range arr {
		if len(arr[i]) > 1 {
			token := strings.Split(arr[i], "/")
			cx := token[1][0]
			flag := false
			for x := range useless {
				if cx == useless[x] {
					//fmt.Println(cx)
					flag = true
				}
			}
			if flag {
				continue
			}
			result = append(result, arr[i])
		}
	}
	return result
}
