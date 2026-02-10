package main

import (
	"fmt"
	"os"
	"the_project/m/v2/routers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"context"
	"time"
	"net/http"
	"strconv"
	"github.com/nats-io/nats.go"
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
var pool *pgxpool.Pool = nil
var nc *nats.Conn = nil

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
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://my-nats:4222"
	}
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Printf("Error connecting to NATS: %v", err)
	} else {
		defer nc.Close()
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/healthz", func(c *gin.Context) {

			ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
			defer cancel()

			if err := pool.Ping(ctx); err != nil {
				c.String(http.StatusServiceUnavailable, "db not ready")
				return
			}

			c.String(http.StatusOK, "ready")
		})


	router.GET("/todos", getTodos)
    router.POST("/todos", createTodo)
	router.PUT("/todos/:id", updateTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started in port %s\n", port)

    router.Run(":" + port)
}


func getTodos(c *gin.Context) {
	var todos []Todo
	rows, err := pool.Query(context.Background(), "SELECT id, title, done FROM todos")
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

	if len(newTodo.Title) > 140 {
        c.JSON(400, gin.H{"error": "Title too long, max 140 characters"})
        return
    }

    _, err := pool.Exec(context.Background(), "INSERT INTO todos (title, done) VALUES ($1, $2)", newTodo.Title, newTodo.Done)
    if err != nil {
        c.JSON(500, gin.H{"error": "Database error"})
        return
    }

	log.Printf("Created todo: %+v\n", newTodo)

	if nc != nil {
		nc.Publish("todo_updates", []byte("A todo was created"))
	}

    c.JSON(201, newTodo)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	_, err = pool.Exec(context.Background(), "UPDATE todos SET done = $1 WHERE id = $2", todo.Done, intID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	log.Printf("Updated todo: %+v\n", todo)

	if nc != nil {
		nc.Publish("todo_updates", []byte("A todo was updated"))
	}

	c.JSON(200, todo)
}

func createhourlyTodo(url string) {

	var err error
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
    
	var newTodo = Todo{
		Title: "Read: " + url,
		Done:  false,
	}	

    _, err = pool.Exec(context.Background(), "INSERT INTO todos (title, done) VALUES ($1, $2)", newTodo.Title, newTodo.Done)
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
