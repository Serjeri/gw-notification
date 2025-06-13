package main

import (
	"context"
	"gw-notification/domain/config"
	mongodb "gw-notification/domain/database"
	"gw-notification/domain/kafka"
	"sync"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.MustLoad()

	log.Info("start application")
	var wg sync.WaitGroup
	workers := 10

	conn, err := mongodb.ConnectDB(cfg.DbURL, cfg.DbName)
	if err != nil {
		log.Info("error connecting to MongoDB: %a", err)
	}
	defer conn.Client().Disconnect(context.Background())

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go kafka.Consume(&wg, i, cfg.KafkaTopic, cfg.Address, cfg.KafkaGroupID, conn)
	}

	wg.Wait()
}
