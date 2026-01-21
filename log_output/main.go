package main

import (
	"crypto/rand"
	"fmt"
	"time"
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

func main() {
	randomString, err := generateRandomString()
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