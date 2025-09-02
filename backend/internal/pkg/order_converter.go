package pkg

import (
	"time"
	"encoding/json"


	"wbts/internal/domain/dto"
	"wbts/internal/domain/entity"
)

type OrderConverter struct{}

func (c *OrderConverter) PaymentDTOToEntity(dto dto.PaymentDTO) entity.Payment {
	return entity.Payment{
		Transaction:  dto.Transaction,
		RequestID:    dto.RequestID,
		Currency:     dto.Currency,
		Provider:     dto.Provider,
		Amount:       dto.Amount,
		PaymentDt:    time.Unix(dto.PaymentDt, 0),
		Bank:         dto.Bank,
		DeliveryCost: dto.DeliveryCost,
		GoodsTotal:   dto.GoodsTotal,
		CustomFee:    dto.CustomFee,
	}
}

func (c *OrderConverter) ItemDTOToEntity(dto dto.ItemDTO) entity.Item {
	return entity.Item{
		ChrtID:      dto.ChrtID,
		TrackNumber: dto.TrackNumber,
		Price:       dto.Price,
		Rid:         dto.Rid,
		Name:        dto.Name,
		Sale:        dto.Sale,
		Size:        dto.Size,
		TotalPrice:  dto.TotalPrice,
		NmID:        dto.NmID,
		Brand:       dto.Brand,
		Status:      dto.Status,
	}
}

func (c *OrderConverter) PaymentEntityToDTO(entity entity.Payment) dto.PaymentDTO {
	return dto.PaymentDTO{
		Transaction:  entity.Transaction,
		RequestID:    entity.RequestID,
		Currency:     entity.Currency,
		Provider:     entity.Provider,
		Amount:       entity.Amount,
		PaymentDt:    entity.PaymentDt.Unix(),
		Bank:         entity.Bank,
		DeliveryCost: entity.DeliveryCost,
		GoodsTotal:   entity.GoodsTotal,
		CustomFee:    entity.CustomFee,
	}
}

func (c *OrderConverter) ItemEntityToDTO(entity entity.Item) dto.ItemDTO {
	return dto.ItemDTO{
		ChrtID:      entity.ChrtID,
		TrackNumber: entity.TrackNumber,
		Price:       entity.Price,
		Rid:         entity.Rid,
		Name:        entity.Name,
		Sale:        entity.Sale,
		Size:        entity.Size,
		TotalPrice:  entity.TotalPrice,
		NmID:        entity.NmID,
		Brand:       entity.Brand,
		Status:      entity.Status,
	}
}

func (c *OrderConverter) OrderDTOToOrderInfo(dto dto.OrderDTO) (entity.OrderInfo, error) {
	payment := c.PaymentDTOToEntity(dto.Payment)

	items := make([]entity.Item, 0, len(dto.Items))
	for _, itemDTO := range dto.Items {
		items = append(items, c.ItemDTOToEntity(itemDTO))
	}

	deliveryJSON, err := json.Marshal(dto.Delivery)
	if err != nil {
		return entity.OrderInfo{}, err
	}

	order := entity.Order{
		OrderUID:          dto.OrderUID,
		TrackNumber:       dto.TrackNumber,
		Entry:             dto.Entry,
		Delivery:          string(deliveryJSON),
		PaymentID:         dto.Payment.Transaction,
		Locale:            dto.Locale,
		InternalSignature: dto.InternalSignature,
		CustomerID:        dto.CustomerID,
		DeliveryService:   dto.DeliveryService,
		Shardkey:          dto.Shardkey,
		SmID:              dto.SmID,
		DateCreated:       dto.DateCreated,
		OofShard:          dto.OofShard,
	}

	return entity.OrderInfo{
		Order:   order,
		Payment: payment,
		Items:   items,
	}, nil
}

func (c *OrderConverter) OrderInfoToOrderDTO(info entity.OrderInfo) (dto.OrderDTO, error) {
	var deliveryDTO dto.DeliveryDTO
	err := json.Unmarshal([]byte(info.Order.Delivery), &deliveryDTO)
	if err != nil {
		return dto.OrderDTO{}, err
	}

	paymentDTO := c.PaymentEntityToDTO(info.Payment)

	itemsDTO := make([]dto.ItemDTO, 0, len(info.Items))
	for _, item := range info.Items {
		itemsDTO = append(itemsDTO, c.ItemEntityToDTO(item))
	}

	return dto.OrderDTO{
		OrderUID:          info.Order.OrderUID,
		TrackNumber:       info.Order.TrackNumber,
		Entry:             info.Order.Entry,
		Delivery:          deliveryDTO,
		Payment:           paymentDTO,
		Items:             itemsDTO,
		Locale:            info.Order.Locale,
		InternalSignature: info.Order.InternalSignature,
		CustomerID:        info.Order.CustomerID,
		DeliveryService:   info.Order.DeliveryService,
		Shardkey:          info.Order.Shardkey,
		SmID:              info.Order.SmID,
		DateCreated:       info.Order.DateCreated,
		OofShard:          info.Order.OofShard,
	}, nil
}
