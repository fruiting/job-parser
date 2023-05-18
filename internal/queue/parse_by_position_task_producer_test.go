package queue

import (
	"errors"
	"fmt"
	"testing"

	"fruiting/job-parser/internal"
	"github.com/golang/mock/gomock"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/suite"
)

type parseByPositionTaskProducerSuite struct {
	suite.Suite

	pool *internal.MockRedisPool

	testErr     error
	payload     *internal.ParseByPositionTask
	payloadJson string

	producer *ParseByPositionTaskProducer
}

func TestParseByPositionTaskProducerSuite(t *testing.T) {
	suite.Run(t, &parseByPositionTaskProducerSuite{})
}

func (s *parseByPositionTaskProducerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.pool = internal.NewMockRedisPool(ctrl)

	s.testErr = errors.New("test err")
	s.payload = &internal.ParseByPositionTask{
		PositionName: "test-position",
	}
	payloadJson, err := easyjson.Marshal(s.payload)
	s.Nil(err)
	s.payloadJson = string(payloadJson)

	s.producer = NewParseByPositionTaskProducer(s.pool)
}

func (s *parseByPositionTaskProducerSuite) TestProduceErr() {
	s.pool.EXPECT().Put(s.payloadJson).Return(s.testErr)

	err := s.producer.Produce(s.payload)

	s.Equal(fmt.Errorf("can't put payload into parse by position task queue: %w", s.testErr), err)
}

func (s *parseByPositionTaskProducerSuite) TestProduceOk() {
	s.pool.EXPECT().Put(s.payloadJson).Return(nil)

	err := s.producer.Produce(s.payload)

	s.Nil(err)
}
