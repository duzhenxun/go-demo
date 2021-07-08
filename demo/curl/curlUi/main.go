package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
	err := ui.Main(func() {
		inputA := ui.NewEntry()
		inputA.SetText("XX活动/交友主播名单")
		inputB := ui.NewEntry()
		inputB.SetText("https://backend.inke.cn/internal/act202106/landlord/user_tab_list?opt=0")
		//生成进度条
		processBar := ui.NewProgressBar()
		processBar.SetValue(0)

		button := ui.NewButton("开始下载")

		//分组
		conT1 := ui.NewGroup("保存地址")
		conT1.SetChild(inputA)

		conT2 := ui.NewGroup("数据下载地址")
		conT2.SetChild(inputB)

		conT3 := ui.NewGroup("")
		conT3.SetChild(processBar)

		conT5 := ui.NewGroup("")
		conT5.SetChild(button)

		box := ui.NewVerticalBox()
		box.Append(conT1, false)
		box.Append(conT2, false)
		box.Append(conT3, false)
		box.Append(conT5, false)

		//创建window窗口。并设置长宽。
		window := ui.NewWindow("活动小助手 V0.0.1（by:DuZhenxun）", 580, 200, true)
		window.SetChild(box)

		//点击事件
		button.OnClicked(func(*ui.Button) {
			//开始处理
			tmpFileName := inputA.Text()
			tmpUrl := inputB.Text()

			if count := strings.Index(tmpUrl, "?"); count == -1 {
				tmpUrl = tmpUrl + "?"
			}
			log.Printf("准备下载数据...  url:%s\n", tmpUrl)
			//time.Sleep(1 * time.Second)
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
			processBar.SetValue(100)
			//greeting.SetText(" 下载完成......\n")

		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func saveData(tmpFileName string, tempMap Res, page int) error {
	delimiter := "/"
	tmpFileName = strings.Trim(tmpFileName, delimiter)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	tmpFileName = dir + delimiter + "data" + delimiter + tmpFileName
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
		err := w.Write(title)
		if err != nil {
			return err
		}
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
		err := w.Write(columnData)
		if err != nil {
			return err
		}
	}
	w.Flush()
	fmt.Printf("============ 数据下载中：%d ============\r", page)
	return nil
}

func mkdir(tmpFileName string, delimiter string) (string, error) {
	count := strings.LastIndex(tmpFileName, delimiter)
	path := ""
	tmpName := tmpFileName
	if count > 0 {
		path = tmpFileName[:count]
		tmpName = tmpFileName[(count + 1):]
	}
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return "", err
	}
	fileName := path + delimiter + tmpName + ".csv"
	return fileName, nil
}
