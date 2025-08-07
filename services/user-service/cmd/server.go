package main

import (
	"log"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	app.Run(cfg)
}
