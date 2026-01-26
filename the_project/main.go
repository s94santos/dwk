package main

import (
	"fmt"
	"os"
	"the_project/m/v2/routers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

const (
	cacheDir  = "cache"
)

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var todos []Todo
var nextID int = 0

func runFrontend() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started in port %s\n", port)

	r := routers.SetRoutes()

	r.Run(":" + port)
}

func runApi() {

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/todos", getTodos)
    router.POST("/todos", createTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started in port %s\n", port)

    router.Run(":" + port)
}


func getTodos(c *gin.Context) {
    c.JSON(200, todos)
}

func createTodo(c *gin.Context) {
    var newTodo Todo
    if err := c.ShouldBindJSON(&newTodo); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request body"})
        return
    }

    newTodo.ID = nextID
    nextID++
    todos = append(todos, newTodo)

    c.JSON(201, newTodo)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: todo-app [api|random]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "backend":
		runApi()
	case "frontend":
		runFrontend()
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}
