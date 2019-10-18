package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main()  {
	r := gin.Default()
	r.StaticFS("/ui", http.Dir("./ui/dist"))

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Second*2)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
