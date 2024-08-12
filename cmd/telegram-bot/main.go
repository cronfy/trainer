package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"

	"oap/trainer/internal/app/useCase/multiplytask"
	bot2 "oap/trainer/internal/telegrambot/bot"
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

	trainerBot := bot2.NewTrainerBot(multiplytask.New())

	opts := []bot.Option{
		bot.WithDefaultHandler(trainerBot.ProcessMessage),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, trainerBot.ProcessStartCommand),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for commands")
	b.Start(ctx)
	fmt.Println("Done")
}
