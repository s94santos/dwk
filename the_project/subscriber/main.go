package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("WEBHOOK_URL is not set, messages will not be forwarded")
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// QueueSubscribe ensures that messages are distributed among the queue group members
	// This allows scaling to multiple replicas without duplicate processing
	_, err = nc.QueueSubscribe("todo_updates", "broadcaster_workers", func(m *nats.Msg) {
		message := string(m.Data)
		log.Printf("Received message: %s", message)

		if webhookURL != "" {
			sendWebhook(webhookURL, message)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Broadcaster started")
	select {} // Block forever
}

func sendWebhook(url, message string) {
	payload := map[string]string{
		"content": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Error sending webhook: %v", err)
		return
	}
	defer resp.Body.Close()
}
