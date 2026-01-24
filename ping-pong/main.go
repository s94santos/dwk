package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"fmt"
	"strconv"
	"io/ioutil"
)

func createFile(cnt int) {
	if _, err := os.Stat("./counter/counter.txt"); err != nil {
		fmt.Println("Creating counter.txt file")

		err = os.WriteFile("./counter/counter.txt", []byte(strconv.Itoa(cnt)), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func writeToFile(cnt int) {

	err := os.WriteFile("./counter/counter.txt", []byte(strconv.Itoa(cnt)), 0644)
	if err != nil {
		panic(err)
	}
}

func getCounter() (int, error) {

	data, err := ioutil.ReadFile("./counter/counter.txt")
    if err != nil {
        panic(err)
    }
	return strconv.Atoi(string(data))
}

func main() {

	createFile(0)

	r := gin.Default()

	r.GET("/pingpong", func(c *gin.Context) {
		cnt, _ := getCounter()
		cnt++
		writeToFile(cnt)
		c.JSON(200, gin.H{
			"counter": cnt,
		})
	})

	r.Run(":8080")

}