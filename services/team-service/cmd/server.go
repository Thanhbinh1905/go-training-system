package main

import (
	"log"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/app"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	app.Run(cfg)
}
