package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

func DecodeFile(w io.Writer, storageName string, nodeDirs []string, exactSize int64) error {

	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}
	shards := make([][]byte, TotalShards)
	for i, node := range nodeDirs {
		path := filepath.Join(node, fmt.Sprintf("%s.shard.%d", storageName, i))
		data, err := os.ReadFile(path)
		// fmt.Printf("Checking shard path: %s\n", path)
		if err == nil {
			shards[i] = data
		} else {
			fmt.Printf("Missing shard %d, will attempt to heal\n", i)
			shards[i] = nil
		}
	}
	err = enc.Reconstruct(shards)
	if err != nil {
		return fmt.Errorf("failed to reconstruct shards: %w", err)
	}
	err = enc.Join(w, shards, int(exactSize))
	if err != nil {
		return fmt.Errorf("failed to join shards: %w", err)
	}
	return nil
}
