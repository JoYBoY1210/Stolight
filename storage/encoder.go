package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

const (
	DataShards   = 4
	ParityShards = 2
	TotalShards  = DataShards + ParityShards
)

func EncodeFile(filePath, outputDir string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open the file: %s", err)
	}
	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %s", err)
	}
	shards, err := enc.Split(data)
	if err != nil {
		return fmt.Errorf("failed to split data into shards: %s", err)
	}
	if err := enc.Encode(shards); err != nil {
		return fmt.Errorf("failed to encode shards: %s", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	fileName := filepath.Base(filePath)
	for i, shard := range shards {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.shard.%d", fileName, i))
		fmt.Printf("Writing shard %d (%d bytes) to %s\n", i, len(shard), outputPath)
		if err := os.WriteFile(outputPath, shard, 0644); err != nil {
			return fmt.Errorf("failed to write shard %d: %s", i, err)
		}
	}
	return nil
}
