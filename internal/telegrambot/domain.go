package telegrambot

import (
	"github.com/go-telegram/bot"

	"oap/trainer/internal/app/domain"
)

type ChatID int64
type State string

const (
	StartState              = "start"
	StateHome               = "home"
	StateWaitingForResponse = "waiting-for-response"
	StateSolvedCorrectly    = "solved-correctly"
	StateSolvedIncorrectly  = "solved-incorrectly"
)

type Session struct {
	State        State
	MultiplyTask *domain.MultiplyTask
}

type Command string

const (
	StartCommand = "start"
)

type EventType int

const (
	MessageEventType EventType = iota
	CommandEventType EventType = iota
)

type Event struct {
	ChatID    ChatID
	B         *bot.Bot
	EventType EventType
	Message   string
	Command   Command
}
