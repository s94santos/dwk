package main

import (
	"fmt"
	"os"
	"the_project/m/v2/routers"
)

const (
	cacheDir  = "cache"
)


func startServer() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started in port %s\n", port)

	r := routers.SetRoutes()

	r.Run(":" + port)
}

func main() {

	_ = os.MkdirAll(cacheDir, 0755)
	startServer()
}
