package main

import (
	"log"

	"github.com/mrsubudei/adv-store-service/internal/app"
	"github.com/mrsubudei/adv-store-service/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}
