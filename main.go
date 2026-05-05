package main

import (
	"context"
	"fmt"

	"slouchdog/tdlib"
)

func main() {
	tdlib := tdlib.Init()
	defer tdlib.Destroy()

	ctx := context.Background()

	go tdlib.Receive(ctx)

	for update := range tdlib.Updates {
		fmt.Println("Received update:", update)
	}
}
