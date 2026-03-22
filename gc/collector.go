package gc

import (
	"context"
	"fmt"
	"time"

	"github.com/joyboy1210/stolight/config"
)

func StartGC(ctx context.Context) {
	if config.Cfg.GCIntervalHours <= 0 {
		fmt.Println("[GC] Garbage Collector disabled in config")
		return
	}
	interval := time.Duration(config.Cfg.GCIntervalHours) * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("[GC] Garbage Collector started, running every %d hours\n", config.Cfg.GCIntervalHours)

	for {
		select {
		case <-ticker.C:
			<-ticker.C
			fmt.Println("[GC] Waking up to clean")

			fmt.Println("[GC] Sweep Done sleeping again.")
		case <-ctx.Done():
			fmt.Println("[GC] Shutting down Garbage Collector")
			return
		}
	}
}
