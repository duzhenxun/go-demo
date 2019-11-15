package main
import (
"github.com/jackdanger/collectlinks"
"net/http"
"fmt"
)

func main() {
	resp, _ := http.Get("https://www.0532888.cn")
	links := collectlinks.All(resp.Body)
	for k,i:=range links{
		fmt.Println(k,i)
	}

}
