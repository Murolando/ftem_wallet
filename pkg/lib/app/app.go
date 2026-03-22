package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/Murolando/ftem_wallet/pkg/config"
)

const (
	ErrorApplication     = "application stopped with error: "
	ErrorFailedToInitApp = "failed to init application: "
)

type App interface {
	Run(ctx context.Context) error
	Init(ctx context.Context) error
}

func RunApplication(config *config.Config, creator func(config *config.Config) App) error {
	var app App

	// add stop signal
	// for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// create app
	if creator == nil {
		app = &BaseApp{}
	} else {
		app = creator(config)
	}

	// init app
	if err := app.Init(ctx); err != nil {
		return errors.New(ErrorFailedToInitApp + err.Error())
	}

	// run app
	if err := app.Run(ctx); err != nil {
		return errors.New(ErrorApplication + err.Error())
	}
	return nil
}
