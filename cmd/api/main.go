package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"wbts/internal/transport/kafka"
)

func main() {
	kafkaCtx, kafkaCtxCancel := context.WithCancel(context.Background())
	c := kafka.NewConsumer(
		[]string{os.Getenv("KAFKA_BROKER")}, 
		os.Getenv("KAFKA_ORDERS_TOPIC"),
		os.Getenv("KAFKA_GROUP_ID"),
	)
	c.Run(kafkaCtx)

	log.Println("Server started")
	http.ListenAndServe(":8081", nil)
	kafkaCtxCancel()
}