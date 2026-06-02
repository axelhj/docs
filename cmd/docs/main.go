package main

import (
	"github.com/axelhj/docs/internal/config"
	"github.com/axelhj/docs/internal/orchestration"
)

func main() {
	cfg := config.Load()

	orchestration.NewApp(cfg)
}
