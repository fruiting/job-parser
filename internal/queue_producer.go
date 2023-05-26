package internal

//go:generate mockgen -source=queue_producer.go -destination=./queue_producer_mock.go -package=internal

type ParseByPositionTaskProducer interface {
	Produce(payload *ParseByPositionTask) error
}
