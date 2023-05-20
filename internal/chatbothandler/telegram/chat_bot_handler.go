package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"fruiting/job-parser/internal"
	"fruiting/job-parser/internal/queue"
	"go.uber.org/zap"
)

const (
	url string = "https://api.telegram.org/bot"
)

type ChatBotHandler struct {
	port                        string
	client                      *http.Client
	apiKey                      string
	parseByPositionTaskProducer *queue.ParseByPositionTaskProducer
	storage                     internal.Storage
	logger                      *zap.Logger
}

func NewChatBotHandler(
	port string,
	client *http.Client,
	apiKey string,
	parseByPositionTaskProducer *queue.ParseByPositionTaskProducer,
	storage internal.Storage,
	logger *zap.Logger,
) *ChatBotHandler {
	return &ChatBotHandler{
		port:                        port,
		client:                      client,
		apiKey:                      apiKey,
		parseByPositionTaskProducer: parseByPositionTaskProducer,
		storage:                     storage,
		logger:                      logger,
	}
}

func (h *ChatBotHandler) ListenAndServe() error {
	return http.ListenAndServe(h.port, http.HandlerFunc(h.handle))
}

func (h *ChatBotHandler) Push(chatId int64, jobsInfo *internal.JobsInfo) error {
	parser := fmt.Sprintf("Parser: %s\n", jobsInfo.Parser)
	position := fmt.Sprintf("Position: %s\n", jobsInfo.PositionToParse)
	minSalary := fmt.Sprintf("MinSalary: %d\n", jobsInfo.MinSalary)
	maxSalary := fmt.Sprintf("MaxSalary: %d\n", jobsInfo.MaxSalary)
	medianSalary := fmt.Sprintf("MedianSalary: %d\n", jobsInfo.MedianSalary)
	popularSkills := "Popular skills:\n"
	for i, skill := range jobsInfo.PopularSkills {
		popularSkills = popularSkills + fmt.Sprintf("%d) %s\n", i+1, skill)
	}
	text := parser + position + minSalary + maxSalary + medianSalary + popularSkills

	err := h.sendMessage(chatId, text)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (h *ChatBotHandler) handle(_ http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("can't ready body", zap.Error(err))

		return
	}

	var request struct {
		Message struct {
			Chat struct {
				Id int64 `json:"id"`
			}
			Text string `json:"text"`
		}
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.logger.Error("can't unmarshal request", zap.Error(err))

		return
	}

	ctxLogger := h.logger.With(zap.Any("request", request.Message.Text))

	command := strings.Split(request.Message.Text, ";")
	if !internal.IsParserInWhiteList(internal.Parser(command[1])) {
		ctxLogger.Warn("requested parser is not in the white list", zap.String("parser", command[1]))
		err = h.sendMessage(request.Message.Chat.Id, "requested parser is not in the white list")
		if err != nil {
			ctxLogger.Error("can't send message to chat bot", zap.Error(err))
		}

		return
	}

	if command[0] == string(internal.ParseJobsInfoChatBotCommand) {
		err = h.parseJobsCommand(command, request.Message.Chat.Id)
	} else if command[0] == string(internal.GetJobsInfoChatBotCommand) {
		err = h.getJobsCommand(req.Context(), command, request.Message.Chat.Id)
	} else {
		err = h.sendMessage(request.Message.Chat.Id, "Invalid command")
	}
	if err != nil {
		ctxLogger.Error("can't handle chat bot request", zap.Error(err))
		err = h.sendMessage(request.Message.Chat.Id, "Whoops! Something went wrong. Check logs.")
		if err != nil {
			ctxLogger.Error("can't send message to chat bot", zap.Error(err))
		}

		return
	}
}

func (h *ChatBotHandler) parseJobsCommand(command []string, chatId int64) error {
	payload := &internal.ParseByPositionTask{
		Parser:       internal.Parser(command[1]),
		ChatId:       chatId,
		PositionName: internal.Name(command[2]),
	}

	err := h.parseByPositionTaskProducer.Produce(payload)
	if err != nil {
		return fmt.Errorf("can't produce message: %w", err)
	}

	err = h.sendMessage(
		chatId,
		fmt.Sprintf("Parsing `%s` is in progress. It may take some time...", command[2]),
	)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (h *ChatBotHandler) getJobsCommand(ctx context.Context, command []string, chatId int64) error {
	jobsInfo, err := h.storage.Get(
		ctx,
		internal.Name(command[1]),
		0, 0, internal.HeadHunterParser)
	if err != nil {
		return fmt.Errorf("can't get jobs info: %w", err)
	}

	err = h.Push(chatId, jobsInfo)
	if err != nil {
		return fmt.Errorf("can't push message: %w", err)
	}

	return nil
}

func (h *ChatBotHandler) sendMessage(chatId int64, text string) error {
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

	res, err := h.client.Post(
		fmt.Sprintf("%s%s/sendMessage", url, h.apiKey),
		"application/json",
		bytes.NewBuffer(msgJson),
	)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("unexpected status: %s", res.Status))
	}

	return nil
}
