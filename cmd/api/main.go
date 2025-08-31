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
	"wbts/internal/transport/rest"
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
		orderService,
		validator,
	)
	go c.Run(context.Background())

	orderHandler := rest.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	mux.HandleFunc("/order/{order_uid}", orderHandler.GetOrderHandler)

	log.Println("Started server on 8081 port")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
