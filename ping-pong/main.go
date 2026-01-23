package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	cnt := 0
	r := gin.Default()

	r.GET("/pingpong", func(c *gin.Context) {
		cnt++
		c.JSON(200, gin.H{
			"counter": cnt,
		})
	})

	r.Run(":8080")

}