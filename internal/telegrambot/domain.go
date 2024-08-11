package telegrambot

type ChatID int64
type State string

const (
	StateHome               = "home"
	StateWaitingForResponse = "waiting-for-response"
	StateSolvedCorrectly    = "solved-correctly"
	StateSolvedIncorrectly  = "solved-incorrectly"
)

type Session struct {
	State State
}
