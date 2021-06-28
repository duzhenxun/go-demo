package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	str := `{
    "data": {
        "header": [
            {
                "title": "序号",
                "key": "num"
            },
            {
                "title": "uid",
                "key": "uid"
            },
            {
                "title": "昵称",
                "key": "nick"
            },
            {
                "title": "积分",
                "key": "score"
            }
        ],
        "column": [
            {
                "uid": "uid",
                "score": "积分",
                "num": "序号",
                "nick": "昵称"
            },
            {
                "uid": 33333,
                "score": 2,
                "num": 2,
                "nick": "爱家族"
            },
            {
                "uid": 22222,
                "score": 2,
                "num": 3,
                "nick": "工"
			}
		]
		}
	}`

	type Header struct {
		Title string `json:"title"`
		Key   string `json:"key"`
	}

	type ResData struct {
		Header []Header                 `json:"header"`
		Column []map[string]interface{} `json:"column"`
	}
	type Res struct {
		Data ResData `json:"data"`
	}
	var tempMap Res
	e := json.Unmarshal([]byte(str), &tempMap)
	fmt.Println(e)

	var title []string
	var titleKey []string
	for _, v := range tempMap.Data.Header {
		title = append(title, fmt.Sprintf("%v", v.Title))
		titleKey = append(titleKey, fmt.Sprintf("%v", v.Key))
	}
	//数据写入表格中
	fileName := "test.csv"
	file, err := os.Open(fileName)
	var w *csv.Writer
	//文件不存在
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
		file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
		//创建一个新的写入文件流
		//写入标题
		w = csv.NewWriter(file)
		w.Write(title)
		w.Write(titleKey)
		w.Flush()
		file.Close()
		fmt.Println("创建新文件")
	}

	fmt.Println("打开文件")
	file, err = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	w = csv.NewWriter(file)
	for _, v := range tempMap.Data.Column {
		var columnData []string
		//这里应根据kk判断一下vv的顺序，保持与表头顺序一致
		for _, tv := range titleKey {
			for kk, vv := range v {
				fmt.Println(kk, vv)
				if tv == kk {
					columnData = append(columnData, fmt.Sprintf("%v", vv))
				}
			}
		}
		w.Write(columnData)
	}
	w.Flush()

}
