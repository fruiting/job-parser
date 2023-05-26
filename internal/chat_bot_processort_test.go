package internal

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type chatBotProcessorSuite struct {
	suite.Suite

	ctx                 context.Context
	testErr             error
	jobsInfo            *JobsInfo
	chatId              int64
	msgTxt              string
	chatBotCommand      *ChatBotCommandInfo
	parseByPositionTask *ParseByPositionTask

	handler                     *MockChatBotHandler
	parseByPositionTaskProducer *MockParseByPositionTaskProducer
	storage                     *MockStorage

	processor *ChatBotProcessor
}

func TestChatBotProcessor(t *testing.T) {
	suite.Run(t, &chatBotProcessorSuite{})
}

func (s *chatBotProcessorSuite) SetupTest() {
	s.ctx = context.Background()
	s.testErr = errors.New("test err")
	s.jobsInfo = &JobsInfo{
		PositionToParse: "golang-developer",
		MinSalary:       100,
		MaxSalary:       200,
		MedianSalary:    200,
		PopularSkills:   []string{"pgsql", "golang", "redis"},
		Parser:          HeadHunterParser,
		Jobs: []*Job{
			{
				PositionName: "golang developer in company 1",
				Link:         "https://hh.ru/1",
				Salary:       100,
				Skills:       []string{"pgsql", "golang", "redis"},
			},
			{
				PositionName: "golang developer in company 2",
				Link:         "https://hh.ru/2",
				Salary:       200,
				Skills:       []string{"golang", "redis"},
			},
			{
				PositionName: "golang developer in company 3",
				Link:         "https://hh.ru/3",
				Salary:       200,
				Skills:       []string{"pgsql", "golang", "redis"},
			},
		},
		Time: time.Now(),
	}
	s.chatId = 800986096
	s.msgTxt = "Parser: hh.ru\n" +
		"Position: golang-developer\n" +
		"MinSalary: 100\n" +
		"MaxSalary: 200\n" +
		"MedianSalary: 200\n" +
		"Popular skills:\n" +
		"1) pgsql\n" +
		"2) golang\n" +
		"3) redis\n"
	s.chatBotCommand = &ChatBotCommandInfo{
		ChatId:   s.chatId,
		Parser:   HeadHunterParser,
		Position: "golang-dev",
		FromYear: 2020,
		ToYear:   2023,
	}
	s.parseByPositionTask = &ParseByPositionTask{
		Parser:       s.chatBotCommand.Parser,
		ChatId:       s.chatBotCommand.ChatId,
		PositionName: s.chatBotCommand.Position,
	}

	ctrl := gomock.NewController(s.T())
	s.handler = NewMockChatBotHandler(ctrl)
	s.parseByPositionTaskProducer = NewMockParseByPositionTaskProducer(ctrl)
	s.storage = NewMockStorage(ctrl)

	s.processor = NewChatBotProcessor(s.handler, s.parseByPositionTaskProducer, s.storage)
}

func (s *chatBotProcessorSuite) TestProcessParseCommandErr() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/parse_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")
	s.handler.EXPECT().ParseCommand(msg).Return(nil, s.testErr)

	err := s.processor.Process(s.ctx, msg)

	s.Equal(fmt.Errorf("can't parse command: %w", s.testErr), err)
}

func (s *chatBotProcessorSuite) TestProcessNotInWhiteListErr() {
	command := s.chatBotCommand
	command.Parser = "some"
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/parse_jobs_info;some;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")

	s.handler.EXPECT().ParseCommand(msg).Return(command, nil)

	err := s.processor.Process(s.ctx, msg)

	s.Equal(errors.New("invalid parser"), err)
}

func (s *chatBotProcessorSuite) TestProcessParseJobsCommandErr() {
	command := s.chatBotCommand
	command.Command = ParseJobsInfoChatBotCommand
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/parse_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")

	s.handler.EXPECT().ParseCommand(msg).Return(command, nil)
	s.parseByPositionTaskProducer.EXPECT().Produce(s.parseByPositionTask).Return(s.testErr)

	err := s.processor.Process(s.ctx, msg)

	s.EqualError(err, "can't process command: can't produce message: test err")
}

func (s *chatBotProcessorSuite) TestProcessParseJobsCommandOk() {
	command := s.chatBotCommand
	command.Command = ParseJobsInfoChatBotCommand
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/parse_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")

	s.handler.EXPECT().ParseCommand(msg).Return(command, nil)
	s.parseByPositionTaskProducer.EXPECT().Produce(s.parseByPositionTask).Return(nil)
	s.handler.
		EXPECT().
		SendMessage(
			s.chatBotCommand.ChatId,
			fmt.Sprintf("Parsing `%s` is in progress. It may take some time...", s.chatBotCommand.Position),
		).
		Return(nil)

	err := s.processor.Process(s.ctx, msg)

	s.Nil(err)
}

func (s *chatBotProcessorSuite) TestProcessGetJobsInfoCommandErr() {
	command := s.chatBotCommand
	command.Command = GetJobsInfoChatBotCommand
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")

	s.handler.EXPECT().ParseCommand(msg).Return(command, nil)
	s.storage.
		EXPECT().
		Get(
			s.ctx,
			s.chatBotCommand.Position,
			s.chatBotCommand.FromYear,
			s.chatBotCommand.ToYear,
			s.chatBotCommand.Parser,
		).
		Return(nil, s.testErr)

	err := s.processor.Process(s.ctx, msg)

	s.EqualError(err, "can't process command: can't get jobs info: test err")
}

func (s *chatBotProcessorSuite) TestProcessGetJobsInfoCommandOk() {
	command := s.chatBotCommand
	command.Command = GetJobsInfoChatBotCommand
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")

	s.handler.EXPECT().ParseCommand(msg).Return(command, nil)
	s.storage.
		EXPECT().
		Get(
			s.ctx,
			s.chatBotCommand.Position,
			s.chatBotCommand.FromYear,
			s.chatBotCommand.ToYear,
			s.chatBotCommand.Parser,
		).
		Return(s.jobsInfo, nil)
	s.handler.EXPECT().SendMessage(s.chatId, s.msgTxt).Return(nil)

	err := s.processor.Process(s.ctx, msg)

	s.Nil(err)
}

func (s *chatBotProcessorSuite) TestProcessInvalidCommandErr() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")

	s.handler.EXPECT().ParseCommand(msg).Return(s.chatBotCommand, nil)
	s.handler.EXPECT().SendMessage(s.chatBotCommand.ChatId, InvalidCommandErr.Error()).Return(s.testErr)

	err := s.processor.Process(s.ctx, msg)

	s.Equal(fmt.Errorf("can't process command: %w", s.testErr), err)
}

func (s *chatBotProcessorSuite) TestProcessInvalidCommandOk() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")

	s.handler.EXPECT().ParseCommand(msg).Return(s.chatBotCommand, nil)
	s.handler.EXPECT().SendMessage(s.chatBotCommand.ChatId, InvalidCommandErr.Error()).Return(nil)

	err := s.processor.Process(s.ctx, msg)

	s.Nil(err)
}

func (s *chatBotProcessorSuite) TestParseJobsCommandProduceErr() {
	s.parseByPositionTaskProducer.EXPECT().Produce(s.parseByPositionTask).Return(s.testErr)

	err := s.processor.parseJobsCommand(s.chatBotCommand)

	s.Equal(fmt.Errorf("can't produce message: %w", s.testErr), err)
}

func (s *chatBotProcessorSuite) TestParseJobsCommandSendMessageErr() {
	s.parseByPositionTaskProducer.EXPECT().Produce(s.parseByPositionTask).Return(nil)
	s.handler.
		EXPECT().
		SendMessage(
			s.chatBotCommand.ChatId,
			fmt.Sprintf("Parsing `%s` is in progress. It may take some time...", s.chatBotCommand.Position),
		).
		Return(s.testErr)

	err := s.processor.parseJobsCommand(s.chatBotCommand)

	s.Equal(fmt.Errorf("can't send message: %w", s.testErr), err)
}

func (s *chatBotProcessorSuite) TestParseJobsCommandOk() {
	s.parseByPositionTaskProducer.EXPECT().Produce(s.parseByPositionTask).Return(nil)
	s.handler.
		EXPECT().
		SendMessage(
			s.chatBotCommand.ChatId,
			fmt.Sprintf("Parsing `%s` is in progress. It may take some time...", s.chatBotCommand.Position),
		).
		Return(nil)

	err := s.processor.parseJobsCommand(s.chatBotCommand)

	s.Nil(err)
}

func (s *chatBotProcessorSuite) TestGetJobsCommandGetErr() {
	s.chatBotCommand.Command = GetJobsInfoChatBotCommand
	s.storage.
		EXPECT().
		Get(
			s.ctx,
			s.chatBotCommand.Position,
			s.chatBotCommand.FromYear,
			s.chatBotCommand.ToYear,
			s.chatBotCommand.Parser,
		).
		Return(nil, s.testErr)

	err := s.processor.getJobsCommand(s.ctx, s.chatBotCommand)

	s.Equal(fmt.Errorf("can't get jobs info: %w", s.testErr), err)
}

func (s *chatBotProcessorSuite) TestGetJobsCommandPushErr() {
	s.chatBotCommand.Command = GetJobsInfoChatBotCommand
	s.storage.
		EXPECT().
		Get(
			s.ctx,
			s.chatBotCommand.Position,
			s.chatBotCommand.FromYear,
			s.chatBotCommand.ToYear,
			s.chatBotCommand.Parser,
		).
		Return(s.jobsInfo, nil)
	s.handler.EXPECT().SendMessage(s.chatId, s.msgTxt).Return(s.testErr)

	err := s.processor.getJobsCommand(s.ctx, s.chatBotCommand)

	s.EqualError(err, fmt.Errorf("can't push message: can't send message: %w", s.testErr).Error())
}

func (s *chatBotProcessorSuite) TestGetJobsCommandOk() {
	s.chatBotCommand.Command = GetJobsInfoChatBotCommand
	s.storage.
		EXPECT().
		Get(
			s.ctx,
			s.chatBotCommand.Position,
			s.chatBotCommand.FromYear,
			s.chatBotCommand.ToYear,
			s.chatBotCommand.Parser,
		).
		Return(s.jobsInfo, nil)
	s.handler.EXPECT().SendMessage(s.chatId, s.msgTxt).Return(nil)

	err := s.processor.getJobsCommand(s.ctx, s.chatBotCommand)

	s.Nil(err)
}

func (s *chatBotProcessorSuite) TestPushErr() {
	s.handler.EXPECT().SendMessage(s.chatId, s.msgTxt).Return(s.testErr)

	err := s.processor.push(s.chatId, s.jobsInfo)

	s.Equal(fmt.Errorf("can't send message: %w", s.testErr), err)
}

func (s *chatBotProcessorSuite) TestPushOk() {
	s.handler.EXPECT().SendMessage(s.chatId, s.msgTxt).Return(nil)

	err := s.processor.push(s.chatId, s.jobsInfo)

	s.Nil(err)
}
