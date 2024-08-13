package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	app "github.com/cronfy/trainer/internal/app/domain"
	tgb "github.com/cronfy/trainer/internal/telegrambot/domain"
	"github.com/cronfy/trainer/internal/telegrambot/sessionstorage"
)

type Bot struct {
	sessionStorage      *sessionstorage.SessionStorage
	multiplyTaskUseCase tgb.MultiplyTaskUseCase
}

func New(multiplyTaskUseCase tgb.MultiplyTaskUseCase, sessionStorage *sessionstorage.SessionStorage) *Bot {
	return &Bot{
		sessionStorage:      sessionStorage,
		multiplyTaskUseCase: multiplyTaskUseCase,
	}
}

func (t *Bot) ProcessCommand(ctx context.Context, b tgb.GoTgBot, update *models.Update) {
	var command tgb.Command

	switch update.Message.Text {
	case "/start":
		command = tgb.StartCommand
	default:
		command = tgb.UnknownCommand
	}

	event := tgb.Event{
		ChatID:    tgb.ChatID(update.Message.Chat.ID),
		EventType: tgb.CommandEventType,
		Message:   update.Message.Text,
		Command:   command,
	}

	t.processEvent(ctx, b, event)
}

func (t *Bot) ProcessMessage(ctx context.Context, b tgb.GoTgBot, update *models.Update) {
	event := tgb.Event{
		ChatID:    tgb.ChatID(update.Message.Chat.ID),
		EventType: tgb.MessageEventType,
		Message:   update.Message.Text,
	}

	t.processEvent(ctx, b, event)
}

func (t *Bot) processEvent(ctx context.Context, b tgb.GoTgBot, event tgb.Event) {
	chatID := event.ChatID
	sess := t.sessionStorage.GetOrCreate(chatID, tgb.Session{State: tgb.StartState})

	if event.EventType == tgb.CommandEventType {
		switch event.Command {
		case tgb.StartCommand:
			sess.State = tgb.StartState
		default:
			t.say(ctx, b, chatID, fmt.Sprintf("Unknown command %s", event.Message))
			return
		}
	}

	nextStep := true
	for nextStep {
		nextStep = false

		switch sess.State {
		case tgb.StartState:
			t.say(ctx, b, chatID, "Hello! I am your trainer. I will give you tasks, reply with an answer.\n\nSend any message when ready.")
			sess.State = tgb.HomeState
		case tgb.HomeState:
			t.giveMultiplyTask(ctx, b, chatID, &sess)
			sess.State = tgb.WaitingForResponseState
		case tgb.WaitingForResponseState:
			if sess.MultiplyTask == nil {
				sess.State = tgb.HomeState
				nextStep = true
				break
			}

			if event.EventType != tgb.MessageEventType {
				break
			}

			if t.verifySolution(*sess.MultiplyTask, event.Message) {
				sess.State = tgb.SolvedCorrectlyState
			} else {
				sess.State = tgb.SolvedIncorrectlyState
			}
			nextStep = true
		case tgb.SolvedCorrectlyState:
			t.say(ctx, b, chatID, "Correct!")
			sess.State = tgb.HomeState
			nextStep = true
		case tgb.SolvedIncorrectlyState:
			t.say(ctx, b, chatID, "Wrong, try again.")
			sess.State = tgb.WaitingForResponseState
		}
	}

	t.sessionStorage.Set(chatID, sess)
}

func (t *Bot) giveMultiplyTask(ctx context.Context, b tgb.GoTgBot, chatID tgb.ChatID, sess *tgb.Session) {
	task := t.multiplyTaskUseCase.Get()
	sess.MultiplyTask = &task
	ops := task.GetOperands()
	t.say(ctx, b, chatID, fmt.Sprintf("%d × %d", ops[0], ops[1]))
}

func (t *Bot) verifySolution(task app.MultiplyTask, response string) bool {
	solution, err := strconv.Atoi(response)
	if err != nil {
		return false
	}
	return t.multiplyTaskUseCase.Solve(task, solution)
}

func (t *Bot) say(ctx context.Context, b tgb.GoTgBot, chatID tgb.ChatID, text string) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatID, Text: text})
}
