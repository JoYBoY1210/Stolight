package gc

import (
	"fmt"
	"os"
	"path/filepath"
)

func SweepStaging() {
	cutoff := GetCutOffTime()
	stagingDir := "./staging"

	contents, err := os.ReadDir(stagingDir)
	if err != nil {
		fmt.Printf("[GC] Error reading staging directory: %v\n", err)
		return
	}
	for _, entry := range contents {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("[GC] Error getting info for %s: %v\n", entry.Name(), err)
			continue
		}
		if info.ModTime().Before(cutoff) {
			path := filepath.Join(stagingDir, entry.Name())
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("[GC] Error removing file %s: %v\n", path, err)
			} else {
				fmt.Printf("[GC] Removed orphaned staging file: %s\n", path)
			}

		}
	}
}
