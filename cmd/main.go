package main

import (
	"context"
	"fmt"
	"gw-notification/domain/database"
	"gw-notification/domain/kafka"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	workers := 3

	conn, err := mongodb.ConnectDB()
	if err != nil {
		fmt.Errorf("error connecting to MongoDB: %a", err)
	}
	defer conn.Client().Disconnect(context.Background())

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go kafka.Consume(&wg, i, "exchange.transaction", "localhost:9092", "exchange.transaction", conn)
	}

	wg.Wait()
}
