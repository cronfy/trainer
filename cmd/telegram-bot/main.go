package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Telegram Bot")

	var err error

	if _, err = os.Stat(".env.local"); err == nil {
		err = godotenv.Load(".env.local", ".env")
	} else {
		err = godotenv.Load(".env")
	}
	if err != nil {
		panic(fmt.Errorf("failed to load env: %w", err))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	trainerBot := NewTrainerBot()

	opts := []bot.Option{
		bot.WithDefaultHandler(trainerBot.processMessage),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for commands")
	b.Start(ctx)
	fmt.Println("Done")
}

type Session struct {
	state State
}

type SessionStorage struct {
	mu       sync.RWMutex
	sessions map[ChatID]Session
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		sessions: make(map[ChatID]Session),
	}
}

func (s *SessionStorage) GetOrCreate(chatId ChatID) Session {
	s.mu.RLock()
	sess, ok := s.sessions[chatId]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		sess = Session{state: StateHome}
		s.sessions[chatId] = sess
		s.mu.Unlock()
	}
	return sess
}

func (s *SessionStorage) Set(chatId ChatID, sess Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatId] = sess
}

type TrainerBot struct {
	sessionStorage *SessionStorage
}

func NewTrainerBot() *TrainerBot {
	return &TrainerBot{
		sessionStorage: NewSessionStorage(),
	}
}

func (t *TrainerBot) giveMultiplyTask(ctx context.Context, b *bot.Bot, chatID ChatID) {
	t.say(ctx, b, chatID, "4 × 12")
}

func (t *TrainerBot) verifySolution(response string) bool {
	return response == "48"
}

func (t *TrainerBot) processMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := ChatID(update.Message.Chat.ID)
	message := update.Message.Text
	sess := t.sessionStorage.GetOrCreate(chatID)

	nextStep := true
	for nextStep {
		nextStep = false

		switch sess.state {
		case StateHome:
			t.giveMultiplyTask(ctx, b, chatID)
			sess.state = StateWaitingForResponse
		case StateWaitingForResponse:
			correct := t.verifySolution(message)
			if correct {
				sess.state = StateSolvedCorrectly
			} else {
				sess.state = StateSolvedIncorrectly
			}
			nextStep = true
		case StateSolvedCorrectly:
			t.say(ctx, b, chatID, "Correct!")
			sess.state = StateHome
			break
		case StateSolvedIncorrectly:
			t.say(ctx, b, chatID, "Wrong, try again.")
			sess.state = StateWaitingForResponse
			break
		}
	}

	t.sessionStorage.Set(chatID, sess)
}

func (t *TrainerBot) say(ctx context.Context, b *bot.Bot, chatID ChatID, text string) {
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: text})

}

type ChatID int64
type State string

const (
	StateHome               = "home"
	StateWaitingForResponse = "waiting-for-response"
	StateSolvedCorrectly    = "solved-correctly"
	StateSolvedIncorrectly  = "solved-incorrectly"
)
