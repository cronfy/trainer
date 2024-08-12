package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	appDomain "github.com/cronfy/trainer/internal/app/domain"
	tb "github.com/cronfy/trainer/internal/telegrambot/domain"
	"github.com/cronfy/trainer/internal/telegrambot/sessionstorage"
)

type TrainerBot struct {
	sessionStorage      *sessionstorage.SessionStorage
	multiplyTaskUseCase tb.MultiplyTaskUseCase
}

func New(multiplyTaskUseCase tb.MultiplyTaskUseCase, sessionStorage *sessionstorage.SessionStorage) *TrainerBot {
	return &TrainerBot{
		sessionStorage:      sessionStorage,
		multiplyTaskUseCase: multiplyTaskUseCase,
	}
}

func (t *TrainerBot) giveMultiplyTask(ctx context.Context, b tb.BotAPI, chatID tb.ChatID, sess *tb.Session) {
	task := t.multiplyTaskUseCase.Get()
	sess.MultiplyTask = &task
	t.say(ctx, b, chatID, fmt.Sprintf("%d × %d", task.Operands[0], task.Operands[1]))
}

func (t *TrainerBot) verifySolution(task appDomain.MultiplyTask, response string) bool {
	solution, err := strconv.Atoi(response)
	if err != nil {
		return false
	}
	return t.multiplyTaskUseCase.Solve(task, solution)
}

func (t *TrainerBot) ProcessCommand(ctx context.Context, b tb.BotAPI, update *models.Update) {
	var command tb.Command

	switch update.Message.Text {
	case "/start":
		command = tb.StartCommand
	default:
		command = tb.UnknownCommand
	}

	event := tb.Event{
		ChatID:    tb.ChatID(update.Message.Chat.ID),
		EventType: tb.CommandEventType,
		Message:   update.Message.Text,
		Command:   command,
	}

	t.processEvent(ctx, b, event)
}

func (t *TrainerBot) ProcessMessage(ctx context.Context, b tb.BotAPI, update *models.Update) {
	event := tb.Event{
		ChatID:    tb.ChatID(update.Message.Chat.ID),
		EventType: tb.MessageEventType,
		Message:   update.Message.Text,
	}

	t.processEvent(ctx, b, event)
}

func (t *TrainerBot) processEvent(ctx context.Context, b tb.BotAPI, event tb.Event) {
	chatID := event.ChatID
	sess := t.sessionStorage.GetOrCreate(chatID)

	if event.EventType == tb.CommandEventType {
		switch event.Command {
		case tb.StartCommand:
			sess.State = tb.StartState
		default:
			t.say(ctx, b, chatID, fmt.Sprintf("Unknown command %s", event.Message))
			return
		}
	}

	nextStep := true
	for nextStep {
		nextStep = false

		switch sess.State {
		case tb.StartState:
			t.say(ctx, b, chatID, "Hello! I am your trainer. I will give you tasks, reply with an answer.\n\nSend any message when ready.")
			sess.State = tb.StateHome
		case tb.StateHome:
			t.giveMultiplyTask(ctx, b, chatID, &sess)
			sess.State = tb.StateWaitingForResponse
		case tb.StateWaitingForResponse:
			if sess.MultiplyTask == nil {
				sess.State = tb.StateHome
				nextStep = true
				break
			}

			if event.EventType != tb.MessageEventType {
				break
			}

			if t.verifySolution(*sess.MultiplyTask, event.Message) {
				sess.State = tb.StateSolvedCorrectly
			} else {
				sess.State = tb.StateSolvedIncorrectly
			}
			nextStep = true
		case tb.StateSolvedCorrectly:
			t.say(ctx, b, chatID, "Correct!")
			sess.State = tb.StateHome
			nextStep = true
		case tb.StateSolvedIncorrectly:
			t.say(ctx, b, chatID, "Wrong, try again.")
			sess.State = tb.StateWaitingForResponse
		}
	}

	t.sessionStorage.Set(chatID, sess)
}

func (t *TrainerBot) say(ctx context.Context, b tb.BotAPI, chatID tb.ChatID, text string) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: text})
}
