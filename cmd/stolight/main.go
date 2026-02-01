package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/db"
)

func main() {
	cfg := config.LoadConfig()

	for _, nodePath := range cfg.StorageNodes {
		if err := os.MkdirAll(nodePath, 0755); err != nil {
			log.Fatalln("Failed to create storage node directory:", err)
		}
	}
	fmt.Println("All nodes initialised")

	Db, err := db.InnitDb(cfg.DBPath)
	if err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}
	fmt.Println("db created successfully")
	db.Mirgrate(Db)
	fmt.Println("tables created successfully")

	fmt.Println("system started successfully")
	select {}

}
