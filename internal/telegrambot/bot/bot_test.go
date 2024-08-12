package bot_test

import (
	"context"
	"regexp"
	"testing"

	gotgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	app "github.com/cronfy/trainer/internal/app/domain"
	"github.com/cronfy/trainer/internal/telegrambot/bot"
	telegrambot "github.com/cronfy/trainer/internal/telegrambot/domain"
	"github.com/cronfy/trainer/internal/telegrambot/domain/mocks"
	"github.com/cronfy/trainer/internal/telegrambot/sessionstorage"
)

type trainerBot interface {
	ProcessCommand(ctx context.Context, b telegrambot.GoTgBot, update *models.Update)
	ProcessMessage(ctx context.Context, b telegrambot.GoTgBot, update *models.Update)
}

type testObjects struct {
	ctx             context.Context
	multiplyUseCase *mocks.MultiplyTaskUseCase
	gotgbot         *mocks.GoTgBot
	sessionStorage  *sessionstorage.SessionStorage
	trainerBot      trainerBot
}

func makeTestObjects(t *testing.T) *testObjects {
	multiplyUseCase := mocks.NewMultiplyTaskUseCase(t)
	sessionStorage := sessionstorage.New()

	return &testObjects{
		ctx:             context.Background(),
		multiplyUseCase: multiplyUseCase,
		gotgbot:         mocks.NewGoTgBot(t),
		sessionStorage:  sessionStorage,
		trainerBot:      bot.New(multiplyUseCase, sessionStorage),
	}
}

// Full multiply task scenario test
func TestTrainerBot_MultiplyScenario(t *testing.T) {
	to := makeTestObjects(t)

	botResponses := make([]string, 0)
	to.gotgbot.EXPECT().SendMessage(mock.Anything, mock.Anything).Run(func(ctx context.Context, params *gotgbot.SendMessageParams) {
		botResponses = append(botResponses, params.Text)
	}).Return(nil, nil)

	var nextMultiplyTask app.MultiplyTask
	var solutionIsCorrect bool
	to.multiplyUseCase.EXPECT().Get().RunAndReturn(func() app.MultiplyTask {
		return nextMultiplyTask
	})
	to.multiplyUseCase.EXPECT().Solve(mock.Anything, mock.Anything).RunAndReturn(func(task app.MultiplyTask, i int) bool {
		return solutionIsCorrect
	})

	t.Run("start", func(t *testing.T) {
		botResponses = botResponses[:0]
		command(t, to, "/start")
		require.Len(t, botResponses, 1)
		require.Contains(t, botResponses[0], "Hello")
	})

	t.Run("send message and get a task", func(t *testing.T) {
		botResponses = botResponses[:0]
		nextMultiplyTask = app.MultiplyTask{Operands: []int{5, 7}}
		message(t, to, "hi")
		require.Len(t, botResponses, 1)
		require.Contains(t, botResponses[0], "5 × 7")
	})

	t.Run("solve the task wrongly", func(t *testing.T) {
		botResponses = botResponses[:0]
		solutionIsCorrect = false
		message(t, to, "37")
		require.Len(t, botResponses, 1)
		require.Contains(t, botResponses[0], "Wrong")
	})

	t.Run("provide not numeric answer", func(t *testing.T) {
		botResponses = botResponses[:0]
		solutionIsCorrect = false
		message(t, to, "i do not know")
		require.Len(t, botResponses, 1)
		require.Contains(t, botResponses[0], "Wrong")
	})

	t.Run("solve the task wrongly again", func(t *testing.T) {
		botResponses = botResponses[:0]
		solutionIsCorrect = false
		message(t, to, "45")
		require.Len(t, botResponses, 1)
		require.Contains(t, botResponses[0], "Wrong")
	})

	t.Run("solve the task correctly", func(t *testing.T) {
		botResponses = botResponses[:0]
		solutionIsCorrect = true
		nextMultiplyTask = app.MultiplyTask{Operands: []int{2, 13}}
		message(t, to, "35")
		require.Len(t, botResponses, 2)
		require.Contains(t, botResponses[0], "Correct")
		require.Contains(t, botResponses[1], "2 × 13")
	})

	t.Run("incorrect command", func(t *testing.T) {
		botResponses = botResponses[:0]
		command(t, to, "/iforgotcommands")
		require.Len(t, botResponses, 1)
		require.Contains(t, botResponses[0], "Unknown command")
	})

	t.Run("restart bot", func(t *testing.T) {
		botResponses = botResponses[:0]
		command(t, to, "/start")
		require.Len(t, botResponses, 1)
		require.Contains(t, botResponses[0], "Hello")
	})
}

// Ensure bot makes correct call to gotgbot api to send response
func TestTrainerBot_ProcessStartCommand_SendMessageCorrectly(t *testing.T) {
	to := makeTestObjects(t)

	type sendMessageCall struct {
		ctx    context.Context
		params gotgbot.SendMessageParams
	}
	sendMessageCalls := make([]sendMessageCall, 0)
	to.gotgbot.EXPECT().SendMessage(mock.Anything, mock.Anything).Run(func(ctx context.Context, params *gotgbot.SendMessageParams) {
		require.NotNil(t, params)
		call := sendMessageCall{ctx: ctx, params: *params}
		sendMessageCalls = append(sendMessageCalls, call)
	}).Return(nil, nil)

	update := models.Update{
		Message: &models.Message{
			Text: "/start",
			Chat: models.Chat{ID: 1122},
		},
	}
	to.trainerBot.ProcessCommand(to.ctx, to.gotgbot, &update)

	require.Len(t, sendMessageCalls, 1)
	assert.Equal(t, sendMessageCalls[0].ctx, to.ctx)
	assert.Equal(t, sendMessageCalls[0].params.ChatID, telegrambot.ChatID(1122))
}

func command(t *testing.T, to *testObjects, command string) {
	t.Helper()

	require.Regexp(t, regexp.MustCompile("^/"), command, "invalid test: command must start with /")

	update := models.Update{
		Message: &models.Message{
			Text: command,
			Chat: models.Chat{},
		},
	}

	to.trainerBot.ProcessCommand(to.ctx, to.gotgbot, &update)
}

func message(t *testing.T, to *testObjects, text string) {
	t.Helper()

	update := models.Update{
		Message: &models.Message{
			Text: text,
			Chat: models.Chat{},
		},
	}

	to.trainerBot.ProcessMessage(to.ctx, to.gotgbot, &update)
}
