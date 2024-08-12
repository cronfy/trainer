package telegrambot

import (
	"github.com/cronfy/trainer/internal/telegrambot/bot"
	"github.com/cronfy/trainer/internal/telegrambot/domain"
	"github.com/cronfy/trainer/internal/telegrambot/sessionstorage"
)

func Build(multiplyTaskUC domain.MultiplyTaskUseCase) *bot.TrainerBot {
	return bot.New(multiplyTaskUC, sessionstorage.New())
}
