package gc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joyboy1210/stolight/config"
)

func SweepNodes() {
	cutoff := GetCutOffTime()

	nodes := config.Cfg.StorageNodes
	for _, node := range nodes {
		entries, err := os.ReadDir(node)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				path := filepath.Join(node, entry.Name())
				err := os.Remove(path)
				if err == nil {
					fmt.Printf("[GC] Swept abandoned .tmp file: %s\n", path)
				}
			}
		}
	}
}
