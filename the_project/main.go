package main

import (
	"fmt"
	"os"
	"the_project/m/v2/routers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/jackc/pgx/v5"
	"log"
	"context"
)

const (
	cacheDir  = "cache"
)

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var nextID int = 0
var conn *pgx.Conn = nil

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
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

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
	var todos []Todo
	rows, err := conn.Query(context.Background(), "SELECT id, title, done FROM todos")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Done)
		if err != nil {
			c.JSON(500, gin.H{"error": "Database error"})
			return
		}
		todos = append(todos, todo)
	}
    c.JSON(200, todos)
}

func createTodo(c *gin.Context) {
    var newTodo Todo	
    if err := c.ShouldBindJSON(&newTodo); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request body"})
        return
    }

    _, err := conn.Exec(context.Background(), "INSERT INTO todos (title, done) VALUES ($1, $2)", newTodo.Title, newTodo.Done)
    if err != nil {
        c.JSON(500, gin.H{"error": "Database error"})
        return
    }

    c.JSON(201, newTodo)
}

func createhourlyTodo(url string) {

	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
    
	var newTodo = Todo{
		Title: "Read: " + url,
		Done:  false,
	}	

    _, err = conn.Exec(context.Background(), "INSERT INTO todos (title, done) VALUES ($1, $2)", newTodo.Title, newTodo.Done)
    if err != nil {
        log.Fatal("Database error:", err)
    }

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: todo-app [backend|frontend|random]")
		os.Exit(1)
	}
	
	switch os.Args[1] {
	case "backend":
		runApi()
	case "frontend":
		runFrontend()
	case "random":
		createhourlyTodo(os.Args[2])
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}
