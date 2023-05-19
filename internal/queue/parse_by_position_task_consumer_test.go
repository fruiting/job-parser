package queue

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"fruiting/job-parser/internal"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type parseByPositionTaskConsumerSuite struct {
	suite.Suite

	ctx           context.Context
	testErr       error
	logs          *observer.ObservedLogs
	generalParser *internal.GeneralParser
	payload       *internal.ParseByPositionTask

	consumer       *internal.MockRedisConsumer
	storage        *internal.MockStorage
	chatBotHandler *internal.MockChatBotHandler
	priceSorter    *internal.MockPriceSorter
	skillsSorter   *internal.MockSkillsSorter
	taskProcessor  *MockTaskProcessor
	jobsParser     *internal.MockJobsParser

	parseByPositionTaskConsumer *ParseByPositionTaskConsumer
}

func TestParseByPositionTaskConsumerSuite(t *testing.T) {
	suite.Run(t, &parseByPositionTaskConsumerSuite{})
}

func (s *parseByPositionTaskConsumerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.ctx = context.Background()
	s.testErr = errors.New("test err")
	core, logs := observer.New(zap.InfoLevel)
	s.logs = logs
	s.payload = &internal.ParseByPositionTask{
		PositionName: "test position",
	}

	s.consumer = internal.NewMockRedisConsumer(ctrl)
	s.storage = internal.NewMockStorage(ctrl)
	s.chatBotHandler = internal.NewMockChatBotHandler(ctrl)
	s.priceSorter = internal.NewMockPriceSorter(ctrl)
	s.skillsSorter = internal.NewMockSkillsSorter(ctrl)
	s.taskProcessor = NewMockTaskProcessor(ctrl)
	s.jobsParser = internal.NewMockJobsParser(ctrl)

	s.generalParser = internal.NewGeneralParser(s.jobsParser)

	s.parseByPositionTaskConsumer = NewParseByPositionTaskConsumer(
		s.consumer,
		s.storage,
		s.chatBotHandler,
		s.priceSorter,
		s.skillsSorter,
		s.generalParser,
		zap.New(core),
	)
}

func (s *parseByPositionTaskConsumerSuite) TestRequeueFailErr() {
	s.taskProcessor.EXPECT().Requeue().Return(s.testErr)
	s.taskProcessor.EXPECT().Fail().Return(s.testErr)

	err := s.parseByPositionTaskConsumer.requeue(s.taskProcessor)

	s.Equal(fmt.Errorf("can't fail task: %w", s.testErr), err)
}

func (s *parseByPositionTaskConsumerSuite) TestRequeueRequeueErr() {
	s.taskProcessor.EXPECT().Requeue().Return(s.testErr)
	s.taskProcessor.EXPECT().Fail().Return(nil)

	err := s.parseByPositionTaskConsumer.requeue(s.taskProcessor)

	s.Equal(fmt.Errorf("can't requeue task: %w", s.testErr), err)
}

func (s *parseByPositionTaskConsumerSuite) TestRequeueOk() {
	s.taskProcessor.EXPECT().Requeue().Return(nil)

	err := s.parseByPositionTaskConsumer.requeue(s.taskProcessor)

	s.Nil(err)
}
