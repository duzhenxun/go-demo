package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Header struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}
type Column struct {
	Uid   int `json:"uid"`
	Score int `json:"score"`
	Num   int `json:"num"`
}

type ResData struct {
	Header []Header `json:"header"`
	Column []Column `json:"column"`
}
type Res struct {
	DmError  int     `json:"dm_error"`
	ErrorMsg string  `json:"error_msg"`
	Data     ResData `json:"data"`
}

func main() {
	url := "https://xxxx.cn/internal/act202106/landlord/user_tab_list?opt=0&page=1"
	resp, err := http.Get(url)
	if err != nil {
		println(err.Error())
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	resp.Body.Close()

	//b := `{"dm_error":0,"error_msg":"\u6210\u529f","data":{"header":[{"title":"\u5e8f\u53f7","key":"num"},{"title":"uid","key":"uid"},{"title":"\u5206\u6570","key":"score"}],"column":[{"uid":712912723,"score":32634,"num":1},{"uid":625130486,"score":14781,"num":2}]}}`

	var res Res
	json.Unmarshal([]byte(b), &res)

	//数据写入表格中
	fileName := "test.csv"
	nfs, _ := os.Create(fileName)
	nfs.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	var data [][]string
	var title []string
	for _, v := range res.Data.Header {
		//fmt.Println(v.Key)
		title = append(title, v.Title)
	}
	data = append(data, title)
	w := csv.NewWriter(nfs)
	w.WriteAll(data)
	w.Flush()
	nfs.Close()

}

func getContent(url string, page string) (res Res, err error) {
	url = url + page
	resp, err := http.Get(url)
	if err != nil {
		println(err.Error())
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(b), &res)

	return
}

func saveCsv(fileName string, content []Column) {
	/*var data [][]string
	for _, v := range content {
		var one []string
		one = append(one, strconv.Itoa(content.Num), strconv.Itoa(v.Uid), strconv.Itoa(v.Score))
		data = append(data, one)
	}
	fp, _ := os.OpenFile(fileName, os.O_RDWR, 0666)
	//fp, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	fp.Seek(0, io.SeekEnd)
	defer fp.Close()
	wData := csv.NewWriter(fp)
	wData.WriteAll(data)
	wData.Flush()*/
}
