package main

import (
	"context"
	"fmt"
	"gw-notification/domain/config"
	mongodb "gw-notification/domain/database"
	"gw-notification/domain/kafka"
	"sync"
)

func main() {
	cfg := config.MustLoad()

	var wg sync.WaitGroup
	workers := 3

	conn, err := mongodb.ConnectDB(cfg.Dburl)
	if err != nil {
		fmt.Errorf("error connecting to MongoDB: %a", err)
	}
	defer conn.Client().Disconnect(context.Background())

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go kafka.Consume(&wg, i, cfg.KafkaTopic, cfg.Address, cfg.KafkaGroupID, conn)
	}

	wg.Wait()
}
