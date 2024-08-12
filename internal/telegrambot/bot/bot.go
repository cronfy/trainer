package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"oap/trainer/internal/app/domain"
	"oap/trainer/internal/telegrambot"
)

type TrainerBot struct {
	sessionStorage      *telegrambot.SessionStorage
	multiplyTaskUseCase multiplyTaskUseCase
}

func NewTrainerBot(multiplyTaskUseCase multiplyTaskUseCase) *TrainerBot {
	return &TrainerBot{
		sessionStorage:      telegrambot.NewSessionStorage(),
		multiplyTaskUseCase: multiplyTaskUseCase,
	}
}

func (t *TrainerBot) giveMultiplyTask(ctx context.Context, b *bot.Bot, chatID telegrambot.ChatID, sess *telegrambot.Session) {
	task := t.multiplyTaskUseCase.Get()
	sess.MultiplyTask = &task
	t.say(ctx, b, chatID, fmt.Sprintf("%d × %d", task.Operands[0], task.Operands[1]))
}

func (t *TrainerBot) verifySolution(task domain.MultiplyTask, response string) bool {
	solution, err := strconv.Atoi(response)
	if err != nil {
		return false
	}
	return t.multiplyTaskUseCase.Solve(task, solution)
}

func (t *TrainerBot) ProcessStartCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	event := telegrambot.Event{
		B:         b,
		ChatID:    telegrambot.ChatID(update.Message.Chat.ID),
		EventType: telegrambot.CommandEventType,
		Command:   telegrambot.StartCommand,
	}

	t.processEvent(ctx, event)
}

func (t *TrainerBot) ProcessMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	event := telegrambot.Event{
		B:         b,
		ChatID:    telegrambot.ChatID(update.Message.Chat.ID),
		EventType: telegrambot.MessageEventType,
		Message:   update.Message.Text,
	}

	t.processEvent(ctx, event)
}

func (t *TrainerBot) processEvent(ctx context.Context, event telegrambot.Event) {
	b := event.B
	chatID := event.ChatID
	sess := t.sessionStorage.GetOrCreate(chatID)

	if event.EventType == telegrambot.CommandEventType {
		switch event.Command {
		case telegrambot.StartCommand:
			sess.State = telegrambot.StartState
		default:
			t.say(ctx, b, chatID, fmt.Sprintf("Unknown command %s", event.Command))
		}
	}

	nextStep := true
	for nextStep {
		nextStep = false

		switch sess.State {
		case telegrambot.StartState:
			t.say(ctx, b, chatID, "Hello! Reply to get a task.")
			sess.State = telegrambot.StateHome
		case telegrambot.StateHome:
			t.giveMultiplyTask(ctx, b, chatID, &sess)
			sess.State = telegrambot.StateWaitingForResponse
		case telegrambot.StateWaitingForResponse:
			if sess.MultiplyTask == nil {
				sess.State = telegrambot.StateHome
				nextStep = true
				break
			}

			if event.EventType != telegrambot.MessageEventType {
				t.say(ctx, b, chatID, "Waiting for solution")
				break
			}

			if t.verifySolution(*sess.MultiplyTask, event.Message) {
				sess.State = telegrambot.StateSolvedCorrectly
			} else {
				sess.State = telegrambot.StateSolvedIncorrectly
			}
			nextStep = true
		case telegrambot.StateSolvedCorrectly:
			t.say(ctx, b, chatID, "Correct!")
			sess.State = telegrambot.StateHome
			nextStep = true
		case telegrambot.StateSolvedIncorrectly:
			t.say(ctx, b, chatID, "Wrong, try again.")
			sess.State = telegrambot.StateWaitingForResponse
		}
	}

	t.sessionStorage.Set(chatID, sess)
}

func (t *TrainerBot) say(ctx context.Context, b *bot.Bot, chatID telegrambot.ChatID, text string) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: text})
}
