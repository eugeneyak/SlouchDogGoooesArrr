package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"slouchdog/bridge"
	"slouchdog/tdlib"
	"syscall"
)

type Typed struct {
	T string `json:"@type"`
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT,
	)
	defer stop()

	bridge, err := bridge.Connect(ctx, "nats://100.104.158.114:4222")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer bridge.Drain()

	td, err := tdlib.Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer td.Destroy()

	updates := td.ReceiveAsync(ctx)

	for update := range updates {
		var s Typed
		err := json.Unmarshal(update, &s)

		if err != nil {
			println(err)
			continue
		}

		log.Println(s.T)

		bridge.EmitUpdate(s.T, update)
	}
}
