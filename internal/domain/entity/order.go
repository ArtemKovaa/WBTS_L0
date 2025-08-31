package entity

import (
    "time"
)

type Payment struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int64
	PaymentDt    time.Time
	Bank         string
	DeliveryCost int64
	GoodsTotal   int64
	CustomFee    int64
}

type Item struct {
	ChrtID      int64
	TrackNumber string
	Price       int64
	Rid         string
	Name        string
	Sale        int8
	Size        string
	TotalPrice  int64
	NmID        int64
	Brand       string
	Status      int
}

type Order struct {
	OrderUID          string
	TrackNumber       string
	Entry             string
	Delivery          string
	PaymentID         string
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	Shardkey          string
	SmID              int64
	DateCreated       time.Time
	OofShard          string
}

type OrderInfo struct {
	Order Order
	Payment Payment
	Items []Item
}
