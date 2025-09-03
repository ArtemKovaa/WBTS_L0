package kafka

import (
    "log"
	"context"
	"encoding/json"

    "github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"

	"wbts/internal/domain/dto"
)

type OrderService interface {
    Save(order dto.OrderDTO)
}

type Consumer struct {
	brokers []string
    topic string
	groupID string
	orderService OrderService
	validator *validator.Validate
}

func NewConsumer(brokers []string, topic string, groupID string, orderService OrderService, validator *validator.Validate) *Consumer {
    return &Consumer{brokers, topic, groupID, orderService, validator}
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
    log.Printf("Start consuming topic: %s", c.topic)
    return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
    log.Printf("Finish consuming topic: %s", c.topic)
    return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        log.Printf(
            "Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
            msg.Topic,
            msg.Partition,
            msg.Offset,
            string(msg.Key),
            string(msg.Value),
        )
		
		var order dto.OrderDTO 
		if err := json.Unmarshal([]byte(msg.Value), &order); err != nil {
			log.Printf("Error parsing message: %v", err)
		}

		if err := c.validator.Struct(order); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
            	log.Printf(
					"Message can't be processed. Field validation error: Field '%s' failed on the '%s' tag\n", 
					err.Field(), 
					err.Tag(),
				)
        	}
		} else {
			c.orderService.Save(order)
		}

        session.MarkMessage(msg, "")
		session.Commit()
    }
    return nil
}

func (c *Consumer) Run(ctx context.Context) {
    consumerGroup, err := sarama.NewConsumerGroup(c.brokers, c.groupID, nil)
    if err != nil {
        log.Fatalf("Error creating consumer group client: %v", err)
    }
    defer consumerGroup.Close()

    for {
        err := consumerGroup.Consume(ctx, []string{c.topic}, c)
        if err != nil {
            log.Fatalf("Error from consumer: %v", err)
        }
        if ctx.Err() != nil {
            return
        }
    }
}
