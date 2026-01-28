package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"context"
	"os"
)

func getCounterFromDB(conn *pgx.Conn) (int, error) {
	var value int
	err := conn.QueryRow(context.Background(), "SELECT value FROM counter WHERE name = 'main'").Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func setCounterInDB(conn *pgx.Conn, value int) error {
	_, err := conn.Exec(context.Background(), "UPDATE counter SET value=$1 WHERE name='main'", value)
	return err
}

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	
	defer conn.Close(context.Background())


	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"healthcheck": "ok",
		})
	})

	r.GET("/pingpong", func(c *gin.Context) {
		cnt, err := getCounterFromDB(conn)
		if err != nil {
			log.Fatal(err)
		}
		cnt++
		err = setCounterInDB(conn, cnt)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"counter": cnt,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		value, err := getCounterFromDB(conn)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"counter": value,
		})
	})

	r.Run(":8080")

}