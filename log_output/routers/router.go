package routers

import (
	"github.com/gin-gonic/gin"
	"time"
	"strconv"
	"io/ioutil"
)

func getCounter() (int, error) {

	data, err := ioutil.ReadFile("./counter/counter.txt")
    if err != nil {
        panic(err)
    }
	return strconv.Atoi(string(data))
}

func SetupRouter(randomString string) *gin.Engine {
	r := gin.Default()

	r.GET("/random", func(c *gin.Context) {
		counter, _ := getCounter()
		c.JSON(200, gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			"value":     randomString,
			"counter":   counter,
		})
	})

	return r
}
