package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joyboy1210/stolight/config"
)

func InitServer(ctx context.Context) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	addr := fmt.Sprintf(":%d", config.Cfg.ServerPort)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		fmt.Printf("Listening on %s\n", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down HTTP server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("HTTP server stopped")
}
