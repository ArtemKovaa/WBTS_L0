package kafka

import (
    "log"
	"context"
    "github.com/IBM/sarama"
)

type Consumer struct {
	brokers []string
    topic string
	groupID string
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
    return &Consumer{brokers, topic, groupID}
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
