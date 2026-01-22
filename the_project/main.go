package main

import (
	"fmt"
	"os"
	"the_project/m/v2/routers"
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
	startServer()
}
