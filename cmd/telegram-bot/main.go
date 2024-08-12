package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"

	"github.com/cronfy/trainer/internal/app/useCase/multiplytask"
	"github.com/cronfy/trainer/internal/telegrambot"
	"github.com/cronfy/trainer/internal/telegrambot/domain"
	"github.com/cronfy/trainer/internal/tools/random"
)

func main() {
	fmt.Println("Starting Telegram Bot")

	loadEnv()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tb := telegrambot.Build(multiplytask.New(random.New()))

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			b := domain.GoTgBot(bot)
			tb.ProcessMessage(ctx, b, update)
		}),
		bot.WithMessageTextHandler(
			"/",
			bot.MatchTypePrefix,
			func(ctx context.Context, bot *bot.Bot, update *models.Update) {
				b := domain.GoTgBot(bot)
				tb.ProcessCommand(ctx, b, update)
			},
		),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for commands")
	b.Start(ctx)
	fmt.Println("Done")
}

func loadEnv() {
	var err error

	if _, err = os.Stat(".env.local"); err == nil {
		err = godotenv.Load(".env.local", ".env")
	} else {
		err = godotenv.Load(".env")
	}
	if err != nil {
		panic(fmt.Errorf("failed to load env: %w", err))
	}
}
