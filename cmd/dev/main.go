package main

import (
	"log"

	"github.com/Murolando/ftem_wallet/internal/cli"
	"github.com/Murolando/ftem_wallet/pkg/config"
	"github.com/Murolando/ftem_wallet/pkg/lib/app"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	if err := app.RunApplication(cfg, func(config *config.Config) app.App {
		return cli.NewCLIController(config)
	}); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
