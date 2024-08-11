package telegrambot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

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

func (t *TrainerBot) ProcessMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := ChatID(update.Message.Chat.ID)
	message := update.Message.Text
	sess := t.sessionStorage.GetOrCreate(chatID)

	nextStep := true
	for nextStep {
		nextStep = false

		switch sess.State {
		case StateHome:
			t.giveMultiplyTask(ctx, b, chatID)
			sess.State = StateWaitingForResponse
		case StateWaitingForResponse:
			correct := t.verifySolution(message)
			if correct {
				sess.State = StateSolvedCorrectly
			} else {
				sess.State = StateSolvedIncorrectly
			}
			nextStep = true
		case StateSolvedCorrectly:
			t.say(ctx, b, chatID, "Correct!")
			sess.State = StateHome
			nextStep = true
		case StateSolvedIncorrectly:
			t.say(ctx, b, chatID, "Wrong, try again.")
			sess.State = StateWaitingForResponse
		}
	}

	t.sessionStorage.Set(chatID, sess)
}

func (t *TrainerBot) say(ctx context.Context, b *bot.Bot, chatID ChatID, text string) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: text})
}
