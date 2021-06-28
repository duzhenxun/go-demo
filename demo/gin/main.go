package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)
var addr = flag.String("addr", ":8080", "register address")
func main()  {
	flag.Parse()
	r := gin.Default()
	r.StaticFS("/ui", http.Dir("./ui/dist"))

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Second*2)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/test",func(c *gin.Context){
		c.JSON(200, map[string]interface{}{"aaa":161526996100001290})

	})
	r.GET("/mysql", mysql)
	r.GET("int", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(*addr) // listen and serve on 0.0.0.0:8080
}



func mysql(c *gin.Context)  {

	db,err := sql.Open("mysql","dy:dy@tcp(127.0.0.1:3306)/dy?charset=utf8");
	if err != nil{
		fmt.Printf("connect mysql fail ! [%s]",err)
	}else{
		fmt.Println("connect to mysql success")
	}

	rows,err := db.Query("select admin_id,admin_name from mac_admin");
	if err != nil{
		fmt.Printf("select fail [%s]",err)
	}

	var mapUser map[string]int
	mapUser = make(map[string]int)

	for rows.Next(){
		var id int
		var username string
		rows.Columns()
		err := rows.Scan(&id,&username)
		if err != nil{
			fmt.Printf("get user info error [%s]",err)
		}
		mapUser[username] = id
	}

	for k,v := range mapUser{
		fmt.Println(k,v);
	}

	db.Close()
}