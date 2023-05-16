package telegram

import (
	"fruiting/job-parser/internal"
	"go.uber.org/zap"
)

const (
	token = "6252436889:AAGFcH73URLj65nQV4v6ZDw7mGI7t6VaRbQ"
	url   = "https://api.telegram.org/bot"
)

type ChatBotHandle struct {
	logger *zap.Logger
}

func NewChatBotHandle(logger *zap.Logger) *ChatBotHandle {
	return &ChatBotHandle{
		logger: logger,
	}
}

func (h *ChatBotHandle) Push(jobsInfo *internal.JobsInfo) error {
	return nil
}
