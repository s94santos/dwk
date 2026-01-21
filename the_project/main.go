package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
)

func startServer() {
	
	port := os.Getenv("PORT")
	fmt.Printf("Server started in port %s\n", port)
	
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Run(":" + port) 
}

func main() {
  startServer()
}

