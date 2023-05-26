package telegram

import (
	"errors"
	"fmt"
	"testing"

	"fruiting/job-parser/internal"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type chatBotHandlerSuite struct {
	suite.Suite

	logs    *observer.ObservedLogs
	testErr error
	apiKey  string

	client  *MockhttpClient
	storage *internal.MockStorage

	handler *ChatBotHandler
}

func TestChatBotHandlerSuite(t *testing.T) {
	suite.Run(t, &chatBotHandlerSuite{})
}

func (s *chatBotHandlerSuite) SetupTest() {
	core, logs := observer.New(zap.InfoLevel)
	s.logs = logs
	s.testErr = errors.New("test err")
	s.apiKey = "test"

	ctrl := gomock.NewController(s.T())
	s.client = NewMockhttpClient(ctrl)
	s.storage = internal.NewMockStorage(ctrl)

	s.handler = NewChatBotHandler(s.client, s.apiKey, s.storage, zap.New(core))
}

func (s *chatBotHandlerSuite) TestParseCommandUnmarshalErr() {
	command, err := s.handler.ParseCommand([]byte(""))

	s.Nil(command)
	s.EqualError(err, "can't unmarshal request: unexpected end of JSON input")
}

func (s *chatBotHandlerSuite) TestParseCommandInvalidCommandErr() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/test\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")
	command, err := s.handler.ParseCommand(msg)

	s.Nil(command)
	s.Equal(internal.InvalidCommandErr, err)
}

func (s *chatBotHandlerSuite) TestParseCommandGetJobsInfoChatBotCommandInvalidCommandErr() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")
	command, err := s.handler.ParseCommand(msg)

	s.Nil(command)
	s.Equal(internal.InvalidCommandErr, err)
}

func (s *chatBotHandlerSuite) TestParseCommandGetJobsInfoChatBotCommandParseFromYearErr() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev;s;s\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")
	command, err := s.handler.ParseCommand(msg)

	s.Nil(command)
	s.EqualError(err, "can't parse from_year: strconv.ParseUint: parsing \"s\": invalid syntax")
}

func (s *chatBotHandlerSuite) TestParseCommandGetJobsInfoChatBotCommandParseToYearErr() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev;2020;s\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")
	command, err := s.handler.ParseCommand(msg)

	s.Nil(command)
	s.EqualError(err, "can't parse to_year: strconv.ParseUint: parsing \"s\": invalid syntax")
}

func (s *chatBotHandlerSuite) TestParseCommandGetJobsInfoChatBotCommandOk() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/get_jobs_info;hh.ru;golang-dev;2020;2023\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")
	command, err := s.handler.ParseCommand(msg)

	s.Equal(&internal.ChatBotCommandInfo{
		ChatId:   328333409,
		Command:  internal.GetJobsInfoChatBotCommand,
		Parser:   internal.HeadHunterParser,
		Position: "golang-dev",
		FromYear: 2020,
		ToYear:   2023,
	}, command)
	s.Nil(err)
}

func (s *chatBotHandlerSuite) TestParseCommandParseJobsInfoChatBotCommandOk() {
	msg := []byte("{\"update_id\":800986096,\n\"message\":{\"message_id\":59,\"from\":{\"id\":328333409,\"is_bot\":false,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"language_code\":\"ru\"},\"chat\":{\"id\":328333409,\"first_name\":\"\\u0420\\u043e\\u043c\\u0430\",\"last_name\":\"\\u0421\\u043f\\u0438\\u0440\\u0438\\u043d\",\"username\":\"romaspirin\",\"type\":\"private\"},\"date\":1685107421,\"text\":\"/parse_jobs_info;hh.ru;golang-dev\",\"entities\":[{\"offset\":0,\"length\":16,\"type\":\"bot_command\"},{\"offset\":17,\"length\":5,\"type\":\"url\"}]}}")
	command, err := s.handler.ParseCommand(msg)

	s.Equal(&internal.ChatBotCommandInfo{
		ChatId:   328333409,
		Command:  internal.ParseJobsInfoChatBotCommand,
		Parser:   internal.HeadHunterParser,
		Position: "golang-dev",
	}, command)
	s.Nil(err)
}

func (s *chatBotHandlerSuite) TestSendMessagePostErr() {
	s.client.
		EXPECT().
		Post(
			fmt.Sprintf("%s%s/%s", tgUrl, s.apiKey, sendMessageUrl),
			[]byte("{\"chat_id\":1234567890,\"text\":\"test\"}"),
		).
		Return(s.testErr)

	err := s.handler.SendMessage(1234567890, "test")

	s.Equal(fmt.Errorf("can't send message: %w", s.testErr), err)
}

func (s *chatBotHandlerSuite) TestSendMessagePosOk() {
	s.client.
		EXPECT().
		Post(
			fmt.Sprintf("%s%s/%s", tgUrl, s.apiKey, sendMessageUrl),
			[]byte("{\"chat_id\":1234567890,\"text\":\"test\"}"),
		).
		Return(nil)

	err := s.handler.SendMessage(1234567890, "test")

	s.Nil(err)
}
