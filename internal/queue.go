package internal

import (
	"github.com/adjust/redismq"
)

//go:generate mockgen -source=queue.go -destination=./queue_mock.go -package=internal
//go:generate easyjson -output_filename=./queue_easyjson.go

// Consumer queue consumer
type Consumer interface {
	Get() (*redismq.Package, error)
}

// easyjson:json Payload
type Payload struct {
	PositionName Name `json:"position_name"`
}
