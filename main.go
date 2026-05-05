package main

import (
	"context"
	"fmt"
	"slouchdog/tdlib"
)

func main() {
	tdlib := tdlib.Init()
	defer tdlib.Destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updates := tdlib.Receive(ctx)

	for update := range updates {
		fmt.Println("Received update:", update)
	}
}
