package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"os"
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

func readFileContent() string {
	content, err := os.ReadFile("./cm/information.txt")
	if err != nil {
		panic(err)
	}
	return string(content)
}

func SetupRouter(randomString string) *gin.Engine {
	r := gin.Default()

	fileContent := readFileContent()

	r.GET("/random", func(c *gin.Context) {
		counter := getCounter()

		c.String(200, "file content: %s env variable: MESSAGE=%s\n%s \nPing / Pongs: %d", fileContent, os.Getenv("MESSAGE"), randomString, counter)
	})

	return r
}
