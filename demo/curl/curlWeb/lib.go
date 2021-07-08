package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func start(fileName string, url string) {
	//开始处理
	var (
		saveFilePath string
		saveFileName string
		tmpUrl       = url
	)

	if count := strings.Index(tmpUrl, "?"); count == -1 {
		tmpUrl = tmpUrl + "?"
	}
	wConn.WriteMessage(1, []byte(fmt.Sprintf("===================  准备下载数据: "+fileName+" =================\n     %s\n", tmpUrl)))
	log.Printf("准备下载数据...  %s\n", tmpUrl)
	//time.Sleep(1 * time.Second)
	//默认拼上page=1
	page := 1
	for true {
		httpUrl := tmpUrl + "&page=" + strconv.Itoa(page)
		resp, err := http.Get(httpUrl)
		if err != nil {
			log.Println("请求接口失败", httpUrl, err)
			wConn.WriteMessage(1, []byte(fmt.Sprintf("请求接口失败,url:%s,err:%v", httpUrl, err)))
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
			if page == 1 {
				log.Println("接口无数据...")
				wConn.WriteMessage(1, []byte("抱歉，你指定的接口无数据...请仔细核对...\n"))
			} else {
				fmt.Println("")
				log.Println("数据下载完成...")
				err := Open(saveFilePath)
				if err != nil {
					fmt.Println(err)
				}
				wConn.WriteMessage(1, []byte(fmt.Sprintf("  恭喜,数据下载完成! %s\n", saveFileName)))
			}

			break
		}

		if saveFilePath, saveFileName, err = saveData(fileName, tempMap, page); err != nil {
			wConn.WriteMessage(1, []byte(fmt.Sprintf("处理异常,err:%v,fileName:%s \n", err, saveFileName)))
			break
		}
		page++
	}
}

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

func saveData(tmpFileName string, tempMap Res, page int) (filePath string, fileName string, err error) {
	delimiter := "/"
	tmpFileName = strings.Trim(tmpFileName, delimiter)
	dir, errAbs := filepath.Abs(filepath.Dir(os.Args[0]))
	if errAbs != nil {
		err = errAbs
		return
	}
	tmpFileName = dir + delimiter + "data" + delimiter + tmpFileName
	filePath, fileName, err = mkdir(tmpFileName, delimiter)
	if err != nil {
		return
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
			return
		}
		file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
		//写入标题
		w = csv.NewWriter(file)
		err = w.Write(title)
		if err != nil {
			return
		}
		w.Write(titleKey)
		w.Flush()
		file.Close()
		time.Sleep(2 * time.Second)
		//log.Printf("创建文件：%s\n", fileName)
		wConn.WriteMessage(1, []byte(fmt.Sprintf("创建文件：%s", fileName)))
	}

	//打开文件，追加数据
	file, err = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		return
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
		err = w.Write(columnData)
		if err != nil {
			return
		}
	}
	wConn.WriteMessage(1, []byte(fmt.Sprintf("  数据下载中....  %d", page)))
	w.Flush()
	fmt.Printf("============ 数据下载中：%d ============\r", page)
	return
}

func mkdir(tmpFileName string, delimiter string) (string, string, error) {
	count := strings.LastIndex(tmpFileName, delimiter)
	path := ""
	tmpName := tmpFileName
	if count > 0 {
		path = tmpFileName[:count]
		tmpName = tmpFileName[(count + 1):]
	}
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return "", "", err
	}
	path += delimiter
	fileName := path + tmpName + ".csv"
	return path, fileName, nil
}

func Open(uri string) error {
	var commands = map[string]string{
		//"windows": "start",
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	//fmt.Println(uri)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start ", uri)
		//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	} else {
		cmd = exec.Command(run, uri)
	}
	return cmd.Start()
}

func hmacSha256(src string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(src))
	shaStr := fmt.Sprintf("%x", h.Sum(nil))
	//shaStr:=hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(shaStr))
}
