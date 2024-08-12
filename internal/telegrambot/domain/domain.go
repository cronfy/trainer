package domain

import (
	"github.com/cronfy/trainer/internal/app/domain"
)

type ChatID int64
type State string

const (
	StartState              State = "start"
	StateHome               State = "home"
	StateWaitingForResponse State = "waiting-for-response"
	StateSolvedCorrectly    State = "solved-correctly"
	StateSolvedIncorrectly  State = "solved-incorrectly"
)

type Session struct {
	State        State
	MultiplyTask *domain.MultiplyTask
}

type Command int

const (
	StartCommand   Command = iota
	UnknownCommand Command = iota
)

type EventType int

const (
	MessageEventType EventType = iota
	CommandEventType EventType = iota
)

type Event struct {
	ChatID    ChatID
	EventType EventType
	Message   string
	Command   Command
}
