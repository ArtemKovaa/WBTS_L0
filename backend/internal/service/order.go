package service

import (
	"log"
	"context"

	"wbts/internal/domain/dto"
	"wbts/internal/domain/entity"
)

type OrderRepo interface {
	GetByUID(ctx context.Context, order_uid string) (*entity.OrderInfo, error) 
	Upsert(ctx context.Context, orderInfo entity.OrderInfo) error
}

type OrderConverter interface {
	OrderDTOToOrderInfo(dto dto.OrderDTO) (entity.OrderInfo, error)
	OrderInfoToOrderDTO(info entity.OrderInfo) (dto.OrderDTO, error)
}

type OrderService struct {
	orderRepo OrderRepo
	orderConverter OrderConverter
}

func NewOrderService(orderRepo OrderRepo, orderConverter OrderConverter) *OrderService {
	return &OrderService {orderRepo, orderConverter}
}

func (s *OrderService) Save(order dto.OrderDTO) {
	orderInfo, err := s.orderConverter.OrderDTOToOrderInfo(order)
	if err != nil {
		log.Printf("Error converting order DTO to Entity: %v", err)
	}

	if err := s.orderRepo.Upsert(context.Background(), orderInfo); err != nil {
		log.Printf("Error saving to DB: %v", err)
	}
}

func (s *OrderService) Get(order_uid string) (dto.OrderDTO, error) {
	orderInfo, err := s.orderRepo.GetByUID(context.Background(), order_uid)
	if err != nil {
		return dto.OrderDTO{}, err
	}

	orderDTO, err := s.orderConverter.OrderInfoToOrderDTO(*orderInfo)
	if err != nil {
		return dto.OrderDTO{}, err
	}

	return orderDTO, nil
}