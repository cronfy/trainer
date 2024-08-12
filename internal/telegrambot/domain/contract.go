package domain

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/cronfy/trainer/internal/app/domain"
)

//go:generate mockery

type BotAPI interface {
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}

type MultiplyTaskUseCase interface {
	Get() domain.MultiplyTask
	Solve(task domain.MultiplyTask, solution int) bool
}
