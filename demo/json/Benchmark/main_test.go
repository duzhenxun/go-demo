package Benchmark

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/buger/jsonparser"
	"testing"
)

var str = `{
        "uid": 214485701,
        "valid_dur": 136192,
        "all_dur": 173268,
        "valid_days": 12,
        "all_days": 59,
        "t61_data": {
            "point_income": "181",
            "live_days": "13"
        }
    }`

func BenchmarkA(b *testing.B) {
	b.ResetTimer()
	data := []byte(str)
	jsonparser.Get(data, "t61_data", "point_income")
	//fmt.Println(string(res))
}

func BenchmarkB(b *testing.B) {
	j, _ := simplejson.NewJson([]byte(str))
	j.Get("t61_data").Get("point_income").String()
	//fmt.Println(s)
}

func BenchmarkC(b *testing.B) {

	var info struct {
		UID       int `json:"uid"`
		ValidDur  int `json:"valid_dur"`  //累计有效开播时长
		AllDur    int `json:"all_dur"`    //累计开播时长
		ValidDays int `json:"valid_days"` //累计有效开播天数
		AllDays   int `json:"all_days"`   //累计开播天数
		T61Data   struct {
			PointIncome string `json:"point_income"` //最近61天 映币收入
			LiveDays    string `json:"live_days"`    //最近61天 开播天数
		} `json:"t61_data"`
	}

	json.Unmarshal([]byte(str), &info)
	//fmt.Println(info.T61Data.PointIncome)
}
