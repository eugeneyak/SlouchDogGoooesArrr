package main

import (
	"context"
	"fmt"
	"os/signal"
	"slouchdog/slouchdog"
	"slouchdog/tdlib"
	"slouchdog/tdlib/update"
	"syscall"
)

var td = tdlib.Init()

func main() {
	defer td.Destroy()

	td.SetLogVerbosityLevel(1)

	ctx, _ := signal.NotifyContext(
		context.Background(), syscall.SIGINT,
	)

	updates := td.Receive(ctx)

	for u := range updates {
		switch v := u.(type) {
		case update.UpdateAuthorizationState:
			slouchdog.Authorize(td, v)

		default:
			fmt.Println("Unhandled update type:", v)
		}
	}
}
