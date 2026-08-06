package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"slouchdog/bridge"
	"slouchdog/tdlib"
	"sync"
	"syscall"

	"github.com/nats-io/nats.go/jetstream"
)

type Typed struct {
	T string `json:"@type"`
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatalln("NATS_URL is empty")
	}

	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT,
	)
	defer stop()

	bridge, err := bridge.Connect(ctx, natsURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer bridge.Drain()

	td, err := tdlib.Init()
	if err != nil {
		log.Fatalln(err)
	}
	td.SetLogVerbosityLevel(tdlib.Fatal)
	defer td.Destroy()

	var wg sync.WaitGroup
	wg.Go(func() { forwardMethods(ctx, td, bridge) })
	wg.Go(func() { forwardUpdates(ctx, td, bridge) })
	wg.Wait()
}

func forwardMethods(ctx context.Context, td *tdlib.TDLib, bridge *bridge.Bridge) {
	err := bridge.Consume(ctx, func(msg jetstream.Msg) {
		td.SendJSON(string(msg.Data()))
		msg.DoubleAck(ctx)
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func forwardUpdates(ctx context.Context, td *tdlib.TDLib, bridge *bridge.Bridge) {
	for update := range td.Updates(ctx) {
		var s Typed

		err := json.Unmarshal([]byte(update), &s)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(s.T)

		bridge.EmitUpdate(s.T, []byte(update))
	}
}
