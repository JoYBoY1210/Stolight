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

	fileIds, err := models.GetAllFileIds()
	if err != nil {
		fmt.Printf("[GC] Failed to retrieve file IDs from DB: %v. Aborting ghost shard sweep.\n", err)
		return
	}

	for _, nodeDir := range config.Cfg.StorageNodes {

		entries, err := os.ReadDir(nodeDir)
		if err != nil {
			fmt.Printf("[GC] Failed to read node directory %s: %v\n", nodeDir, err)
			continue
		}

		for _, entry := range entries {

			if entry.IsDir() {
				continue
			}

			name := entry.Name()

			if strings.HasSuffix(name, ".tmp") {
				continue
			}

			idx := strings.Index(name, ".shard.")
			if idx == -1 {
				continue
			}
			if idx == 0 {
				continue
			}

			fileID := name[:idx]

			if fileIds[fileID] {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			
			if !info.ModTime().Before(cutoff) {
				continue
			}

			fullPath := filepath.Join(nodeDir, name)

			if err := os.Remove(fullPath); err != nil {
				fmt.Printf("[GC] Failed to delete ghost shard %s: %v\n", name, err)
				continue
			}

			fmt.Printf("[GC] Exorcised ghost shard: %s\n", name)
		}
	}
}
