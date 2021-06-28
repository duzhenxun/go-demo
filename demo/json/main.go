package main

import (
	"fmt"
	"github.com/buger/jsonparser"
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

func main() {

	data := []byte(str)
	b, t, o, e := jsonparser.Get(data, "t61_data")
	fmt.Println(string(b), t, o, e)

}
