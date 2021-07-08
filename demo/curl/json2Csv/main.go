package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//str := `{
//    "data": {
//        "header": [
//            {
//                "title": "序号",
//                "key": "num"
//            },
//            {
//                "title": "uid",
//                "key": "uid"
//            },
//            {
//                "title": "昵称",
//                "key": "nick"
//            },
//            {
//                "title": "积分",
//                "key": "score"
//            }
//        ],
//        "column": [
//            {
//                "uid": "uid",
//                "score": "积分",
//                "num": "序号",
//                "nick": "昵称"
//            },
//            {
//                "uid": 33333,
//                "score": 2,
//                "num": 2,
//                "nick": "爱家族"
//            },
//            {
//                "uid": 22222,
//                "score": 2,
//                "num": 3,
//                "nick": "工"
//			}
//		]
//		}
//	}`

//	./main -url="https://backend.inke.cn/internal/act202106/landlord/user_tab_list?opt=0" -file_name=斗地主/主播榜个人

var fileName = flag.String("file_name", "xxxx活动/xx数据", "要保存的文件名，xxxx活动/xx数据")
var url = flag.String("url", "https://backend.inke.cn/internal/act202106/landlord/user_tab_list?opt=0", "输入后台地址")

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

func main() {
	flag.Parse()
	tmpFileName := *fileName
	tmpUrl := *url
	if len := strings.Index(tmpUrl, "?"); len == -1 {
		tmpUrl = tmpUrl + "?"
	}
	//exec.Command(`open`, `https://xs25.cn`).Start()

	log.Printf("准备下载数据...  url:%s\n", tmpUrl)
	time.Sleep(2 * time.Second)
	//默认拼上page=1
	page := 1
	for true {
		httpUrl := tmpUrl + "&page=" + strconv.Itoa(page)
		resp, err := http.Get(httpUrl)
		if err != nil {
			log.Println("请求接口失败", httpUrl, err)
			break
		}
		b, ReadErr := ioutil.ReadAll(resp.Body)
		if ReadErr != nil {
			log.Println("读取数据失败", httpUrl, ReadErr)
			break
		}
		str := string(b)
		d := json.NewDecoder(strings.NewReader(str))
		d.UseNumber()

		var tempMap Res

		if jsonErr := d.Decode(&tempMap); jsonErr != nil {
			log.Println("json解析失败", httpUrl, ReadErr, string(str))

		}

		//查看数据
		if len(tempMap.Data.Column) == 0 {
			fmt.Println("")
			log.Println("数据下载完成...")
			break
		}
		saveData(tmpFileName, tempMap, page)
		page++
	}

}

func saveData(tmpFileName string, tempMap Res, page int) error {
	delimiter := "/"
	tmpFileName = strings.Trim(tmpFileName, delimiter)
	tmpFileName = "data" + delimiter + tmpFileName
	fileName, err := mkdir(tmpFileName, delimiter)
	if err != nil {
		return err
	}

	var (
		title    []string
		titleKey []string
	)
	for _, v := range tempMap.Data.Header {
		title = append(title, fmt.Sprintf("%v", v.Title))
		titleKey = append(titleKey, fmt.Sprintf("%v", v.Key))
	}

	var file *os.File
	var w *csv.Writer
	if page == 1 {
		//每次重新创建新文件
		file, err = os.Create(fileName)
		if err != nil {
			return err
		}
		file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
		//写入标题
		w = csv.NewWriter(file)
		w.Write(title)
		w.Write(titleKey)
		w.Flush()
		file.Close()
		log.Printf("创建文件：%s\n", fileName)
		time.Sleep(2 * time.Second)
	}

	//打开文件，追加数据
	file, err = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	w = csv.NewWriter(file)
	for _, v := range tempMap.Data.Column {
		var columnData []string
		//这里应根据kk判断一下vv的顺序，保持与表头顺序一致
		for _, val := range titleKey {
			for kk, vv := range v {
				if val == kk {
					columnData = append(columnData, fmt.Sprintf("%v", vv))
				}
			}
		}
		//添加一条数据
		w.Write(columnData)
	}
	w.Flush()
	fmt.Printf("============ 数据下载中：%d ============\r", page)
	return nil
}

func mkdir(tmpFileName string, delimiter string) (string, error) {
	len := strings.LastIndex(tmpFileName, delimiter)
	path := ""
	tmpName := tmpFileName
	if len > 0 {
		path = tmpFileName[:len]
		tmpName = tmpFileName[(len + 1):]
	}
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return "", err
	}
	fileName := path + delimiter + tmpName + ".csv"
	return fileName, nil
}
