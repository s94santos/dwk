package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"context"
	"os"
	"time"
	"net/http"
)

type App struct {
	db *pgxpool.Pool
	router *gin.Engine
}

func (app *App) getCounterFromDB() (int, error) {
	var value int
	err := app.db.QueryRow(context.Background(), "SELECT value FROM counter WHERE name = 'main'").Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (app *App) setCounterInDB(value int) error {
	_, err := app.db.Exec(context.Background(), "UPDATE counter SET value=$1 WHERE name='main'", value)
	return err
}

func main() {

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	r := gin.Default()

	app := App{
		db: pool,
		router: r,
	}

	r.GET("/healthz", func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
		defer cancel()

		if err := app.db.Ping(ctx); err != nil {
			c.String(http.StatusServiceUnavailable, "db not ready")
			return
		}

		c.String(http.StatusOK, "ready")
	})

	r.GET("/", func(c *gin.Context) {
		cnt, err := app.getCounterFromDB()
		if err != nil {
			log.Fatal(err)
		}
		cnt++
		err = app.setCounterInDB(cnt)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"counter": cnt,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		value, err := app.getCounterFromDB()
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"counter": value,
		})
	})

	r.Run(":8080")

}