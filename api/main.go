package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joyboy1210/stolight/config"
)

func InitServer() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := mux
	addr := fmt.Sprintf(":%d", config.Cfg.ServerPort)
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Printf("Listening on %s\n", addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server could not be started: %v\n", err)
	}
}
