package queue

import (
	"fmt"

	"fruiting/job-parser/internal"
	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	processors []*internal.ParsingProcessor
}

func NewConsumerHandler(processors []*internal.ParsingProcessor) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		processors: processors,
	}
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}

	return nil
}
