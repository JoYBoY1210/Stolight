package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/db"
	"github.com/joyboy1210/stolight/handlers"
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

	uploadHandler := &handlers.UploadHandler{
		Db:     Db,
		Config: cfg,
	}
	downloadHandler := &handlers.DownloadRequest{
		Db: Db,
	}
	http.HandleFunc("/upload", uploadHandler.UploadHandlerAPI)
	http.HandleFunc("/download/", downloadHandler.DownloadHandler)
	fmt.Printf("Listening on :%d\n", cfg.ServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), nil); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}

}
