package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/klauspost/reedsolomon"
)

func DecodeFile(shardDir, fileName, outputPath string) error {
	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return fmt.Errorf("could not create the encoder: %s", err)
	}
	shards := make([][]byte, TotalShards)
	for i := 0; i < TotalShards; i++ {
		shardPath := filepath.Join(shardDir, fmt.Sprintf("%s.shard.%d", fileName, i))
		shardData, err := os.ReadFile(shardPath)
		if err != nil {
			return fmt.Errorf("could not read shard %d: %s", i, err)
		} else {
			shards[i] = shardData
			fmt.Printf("found shard number %d\n", i)
		}

	}
	err = enc.Reconstruct(shards)
	if err != nil {
		return fmt.Errorf("could not reconstruct the data: %s", err)
	}
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %s", err)
	}
	defer outFile.Close()
	for i := 0; i < DataShards; i++ {
		outFile.Write(shards[i])
	}
	return nil
}
