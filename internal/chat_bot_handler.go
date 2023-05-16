package internal

//go:generate mockgen -source=chat_bot_handler.go -destination=./chat_bot_handler_mock.go -package=internal

// ChatBotHandler handles with chatbot
type ChatBotHandler interface {
	Push(jobsInfo *JobsInfo) error
}
