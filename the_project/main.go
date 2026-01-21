package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
)

func startServer() {
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started in port %s\n", port)
	
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Run(":" + port) 
}

func main() {
  startServer()
}

