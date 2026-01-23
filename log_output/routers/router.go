package routers

import (
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter(randomString string) *gin.Engine {
	r := gin.Default()

	r.GET("/random", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			"value":     randomString,
		})
	})

	return r
}
