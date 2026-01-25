package routers

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"encoding/json"
)

type PingResponse struct {
	Counter int `json:"counter"`
}

func getCounter() int {
	resp, err := http.Get("http://ping-pong-svc:3000/ping")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data PingResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	return data.Counter
}

func SetupRouter(randomString string) *gin.Engine {
	r := gin.Default()

	r.GET("/random", func(c *gin.Context) {
		counter := getCounter()
		c.JSON(200, gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			"value":     randomString,
			"counter":   counter,
		})
	})

	return r
}
