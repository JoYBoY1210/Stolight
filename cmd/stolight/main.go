package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joyboy1210/stolight/api"
	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/db"
	"github.com/joyboy1210/stolight/gc"
	"github.com/joyboy1210/stolight/models"
	"github.com/joyboy1210/stolight/queue"
	"github.com/joyboy1210/stolight/utils"
)

func main() {

	for _, nodePath := range config.Cfg.StorageNodes {
		if err := os.MkdirAll(nodePath, 0755); err != nil {
			log.Fatalln("Failed to create storage node directory:", err)
		}
	}
	fmt.Println("All nodes initialised")

	Db, err := db.InnitDb(config.Cfg.DBPath)
	if err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	models.SetDB(Db)
	fmt.Println("db created successfully")
	db.Mirgrate(Db)
	utils.CheckOnStart(Db)
	q := queue.InitQueue(ctx, 100)
	go gc.StartGC(ctx)
	api.InitServer(ctx)
	<-ctx.Done()
	q.Close()

}
