package gc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/models"
)

func SweepGhostShards() {
	cutoff := GetCutOffTime()

	for _, nodeDir := range config.Cfg.StorageNodes {
		entries, err := os.ReadDir(nodeDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()

			if !strings.Contains(name, ".shard.") || strings.HasSuffix(name, ".tmp") {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.ModTime().Before(cutoff) {
				parts := strings.Split(name, ".shard.")
				if len(parts) < 2 {
					continue
				}
				fileID := parts[0]

				_, err := models.GetFileByID(fileID)
				if err != nil {
					fullPath := filepath.Join(nodeDir, name)
					err := os.Remove(fullPath)
					if err == nil {
						fmt.Printf("[GC] Exorcised ghost shard: %s\n", name)
					} else {
						fmt.Printf("[GC] Failed to delete ghost %s: %v\n", name, err)
					}
				}
			}
		}
	}
}
