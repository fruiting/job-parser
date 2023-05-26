package internal

//go:generate mockgen -source=chat_bot_processor.go -destination=./chat_bot_processor_mock.go -package=internal

import (
	"context"
	"errors"
	"fmt"
)

// ChatBotHandler handles with chat bot
type ChatBotHandler interface {
	ParseCommand(bodyRequest []byte) (*ChatBotCommandInfo, error)
	SendMessage(chatId int64, text string) error
}

const (
	ParseJobsInfoChatBotCommand ChatBotCommand = "/parse_jobs_info"
	GetJobsInfoChatBotCommand   ChatBotCommand = "/get_jobs_info"
)

type ChatBotCommand string

type ChatBotCommandInfo struct {
	ChatId   int64
	Command  ChatBotCommand
	Parser   Parser
	Position Name
	FromYear uint16
	ToYear   uint16
}

type ChatBotProcessor struct {
	handler                     ChatBotHandler
	parseByPositionTaskProducer ParseByPositionTaskProducer
	storage                     Storage
}

func NewChatBotProcessor(
	handler ChatBotHandler,
	parseByPositionTaskProducer ParseByPositionTaskProducer,
	storage Storage,
) *ChatBotProcessor {
	return &ChatBotProcessor{
		handler:                     handler,
		parseByPositionTaskProducer: parseByPositionTaskProducer,
		storage:                     storage,
	}
}

func (p *ChatBotProcessor) Process(ctx context.Context, body []byte) error {
	chatBotCommand, err := p.handler.ParseCommand(body)
	if err != nil {
		return fmt.Errorf("can't parse command: %w", err)
	}

	if !IsParserInWhiteList(chatBotCommand.Parser) {
		return errors.New("invalid parser")
	}

	if chatBotCommand.Command == ParseJobsInfoChatBotCommand {
		err = p.parseJobsCommand(chatBotCommand)
	} else if chatBotCommand.Command == GetJobsInfoChatBotCommand {
		err = p.getJobsCommand(ctx, chatBotCommand)
	} else {
		err = p.handler.SendMessage(chatBotCommand.ChatId, InvalidCommandErr.Error())
	}
	if err != nil {
		return fmt.Errorf("can't process command: %w", err)
	}

	return nil
}

func (p *ChatBotProcessor) parseJobsCommand(chatBotCommand *ChatBotCommandInfo) error {
	payload := &ParseByPositionTask{
		Parser:       chatBotCommand.Parser,
		ChatId:       chatBotCommand.ChatId,
		PositionName: chatBotCommand.Position,
	}

	err := p.parseByPositionTaskProducer.Produce(payload)
	if err != nil {
		return fmt.Errorf("can't produce message: %w", err)
	}

	err = p.handler.SendMessage(
		chatBotCommand.ChatId,
		fmt.Sprintf("Parsing `%s` is in progress. It may take some time...", chatBotCommand.Position),
	)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (p *ChatBotProcessor) getJobsCommand(ctx context.Context, chatBotCommand *ChatBotCommandInfo) error {
	jobsInfo, err := p.storage.Get(
		ctx,
		chatBotCommand.Position,
		chatBotCommand.FromYear,
		chatBotCommand.ToYear,
		chatBotCommand.Parser,
	)
	if err != nil {
		return fmt.Errorf("can't get jobs info: %w", err)
	}

	err = p.push(chatBotCommand.ChatId, jobsInfo)
	if err != nil {
		return fmt.Errorf("can't push message: %w", err)
	}

	return nil
}

func (p *ChatBotProcessor) push(chatId int64, jobsInfo *JobsInfo) error {
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

	err := p.handler.SendMessage(chatId, text)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}
