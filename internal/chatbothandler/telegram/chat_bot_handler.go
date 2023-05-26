package telegram

//go:generate mockgen -source=chat_bot_handler.go -destination=./chat_bot_handler_mock.go -package=telegram

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"fruiting/job-parser/internal"
	"go.uber.org/zap"
)

type httpClient interface {
	Post(url string, body []byte) error
}

const (
	tgUrl          string = "https://api.telegram.org/bot"
	sendMessageUrl string = "sendMessage"
)

type ChatBotHandler struct {
	client  httpClient
	apiKey  string
	storage internal.Storage
	logger  *zap.Logger
}

func NewChatBotHandler(
	client httpClient,
	apiKey string,
	storage internal.Storage,
	logger *zap.Logger,
) *ChatBotHandler {
	return &ChatBotHandler{
		client:  client,
		apiKey:  apiKey,
		storage: storage,
		logger:  logger,
	}
}

func (h *ChatBotHandler) ParseCommand(bodyRequest []byte) (*internal.ChatBotCommandInfo, error) {
	var request struct {
		Message struct {
			Chat struct {
				Id int64 `json:"id"`
			}
			Text string `json:"text"`
		}
	}

	err := json.Unmarshal(bodyRequest, &request)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal request: %w", err)
	}

	command := strings.Split(request.Message.Text, ";")
	if len(command) < 3 {
		return nil, internal.InvalidCommandErr
	}

	commandInfo := &internal.ChatBotCommandInfo{
		ChatId:   request.Message.Chat.Id,
		Command:  internal.ChatBotCommand(command[0]),
		Parser:   internal.Parser(command[1]),
		Position: internal.Name(command[2]),
	}
	if command[0] == string(internal.GetJobsInfoChatBotCommand) {
		if len(command) < 5 {
			return nil, internal.InvalidCommandErr
		}

		fromYear, err := strconv.ParseUint(command[3], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("can't parse from_year: %w", err)
		}

		toYear, err := strconv.ParseUint(command[4], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("can't parse to_year: %w", err)
		}

		commandInfo.FromYear = uint16(fromYear)
		commandInfo.ToYear = uint16(toYear)
	}

	return commandInfo, nil
}

func (h *ChatBotHandler) SendMessage(chatId int64, text string) error {
	type sendMessageBody struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}

	msg := &sendMessageBody{
		ChatId: chatId,
		Text:   text,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("can't marshal msg: %w", err)
	}

	err = h.client.Post(fmt.Sprintf("%s%s/%s", tgUrl, h.apiKey, sendMessageUrl), msgJson)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}
