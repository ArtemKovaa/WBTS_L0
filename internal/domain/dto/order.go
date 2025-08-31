package dto

import (
	"time"
)

type DeliveryDTO struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required,min=3,max=32"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type PaymentDTO struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int64  `json:"amount" validate:"gte=0"`
	PaymentDt    int64  `json:"payment_dt" validate:"gte=0"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int64  `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int64  `json:"goods_total" validate:"gte=0"`
	CustomFee    int64  `json:"custom_fee" validate:"gte=0"`
}

type ItemDTO struct {
	ChrtID      int64  `json:"chrt_id" validate:"gt=0"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price" validate:"gte=0"`
	Rid         string `json:"rid"`
	Name        string `json:"name" validate:"required"`
	Sale        int8   `json:"sale" validate:"gte=0,lte=100"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int64  `json:"total_price" validate:"gt=0"`
	NmID        int64  `json:"nm_id" validate:"gt=0"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

type OrderDTO struct {
	OrderUID          string      `json:"order_uid" validate:"required"`
	TrackNumber       string      `json:"track_number"`
	Entry             string      `json:"entry" validate:"required"`
	Delivery          DeliveryDTO `json:"delivery" validate:"required"`
	Payment           PaymentDTO  `json:"payment" validate:"required"`
	Items             []ItemDTO   `json:"items" validate:"required"`
	Locale            string      `json:"locale" validate:"required"`
	InternalSignature string      `json:"internal_signature" validate:"required"`
	CustomerID        string      `json:"customer_id" validate:"required"`
	DeliveryService   string      `json:"delivery_service" validate:"required"`
	Shardkey          string      `json:"shardkey" validate:"required"`
	SmID              int64       `json:"sm_id" validate:"gt=0"`
	DateCreated       time.Time   `json:"date_created" validate:"required"`
	OofShard          string      `json:"oof_shard" validate:"required"`
}
