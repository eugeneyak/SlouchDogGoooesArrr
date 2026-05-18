package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"slouchdog/slouchdog"
	"slouchdog/tdlib"
	"slouchdog/tdlib/log"
	"slouchdog/tdlib/update"
)

var wg sync.WaitGroup
var td *tdlib.TDLib

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT,
	)

	defer stop()

	wg.Add(2)

	go starttdlib(ctx)
	go startweb(ctx)

	wg.Wait()
}

func starttdlib(ctx context.Context) {
	td := tdlib.Init()
	defer td.Destroy()

	td.SetLogVerbosityLevel(log.Error)

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

func startweb(ctx context.Context) {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Static("./dogface/dist"))

	sc := echo.StartConfig{
		Address:         ":1323",
		GracefulTimeout: 5 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
