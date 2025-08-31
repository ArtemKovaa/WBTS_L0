package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"

	"wbts/internal/pkg"
	"wbts/internal/service"
	"wbts/internal/storage"
	"wbts/internal/transport/kafka"
)

func main() {
	pgPool := storage.Setup(context.Background())
	defer pgPool.Close()
	orderRepo := storage.NewOrderRepo(pgPool)
	orderConverter := &pkg.OrderConverter{}
	orderService := service.NewOrderService(orderRepo, orderConverter)
	validator := validator.New()

	c := kafka.NewConsumer(
		[]string{os.Getenv("KAFKA_BROKER")},
		os.Getenv("KAFKA_ORDERS_TOPIC"),
		os.Getenv("KAFKA_GROUP_ID"),
		context.Background(),
		orderService,
		validator,
	)
	c.Run(context.Background())

	log.Println("Server started")
	http.ListenAndServe(":8081", nil)
}
