package main

import (
	"crypto/rand"
	"fmt"
	"time"
	"log_module/m/v2/routers"
	"os"
)

func generateRandomString() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// Format as UUID v4 style
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}

func createFile() {
	if _, err := os.Stat("./data/data.txt"); err != nil {
		fmt.Println("Creating data.txt file")
		randomString, err := generateRandomString()
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("./data/data.txt", []byte(randomString), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func runRandom() {

	randomString, err := os.ReadFile("./data/data.txt")
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		t := <-ticker.C
		fmt.Printf("%s: %s\n", t.UTC().Format(time.RFC3339Nano), randomString)
	}
}

func runApi() {
	randomString, err := os.ReadFile("./data/data.txt")
	if err != nil {
		panic(err)
	}

	router := routers.SetupRouter(string(randomString))
	router.Run(":9090")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: todo-app [api|random]")
		os.Exit(1)
	}

	createFile()

	switch os.Args[1] {
	case "api":
		runApi()
	case "random":
		runRandom()
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}
