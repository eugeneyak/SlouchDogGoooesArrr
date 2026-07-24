package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
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

	// nc, err := bridge.Connect("nats://100.104.158.114:4222")
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)
	// }
	// defer nc.Drain()

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

		// 	msg := nats.NewMsg("tdlib.update")

		// 	msg.Header.Add("type", s.T)
		// 	msg.Data = []byte(update)

		// 	nc.PublishMsg(msg)
	}
}
