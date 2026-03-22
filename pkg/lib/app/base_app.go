package app

import (
	"context"
	"fmt"
)

type BaseApp struct {
}

func (a *BaseApp) Init(ctx context.Context) error {
	fmt.Println("Init BaseApp")
	return nil
}

func (a *BaseApp) Run(ctx context.Context) error {
	fmt.Println("Run BaseApp")
	return nil
}
