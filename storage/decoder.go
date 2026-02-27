package storage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joyboy1210/stolight/models"
	"github.com/klauspost/reedsolomon"
)

func DecodeFile(w io.Writer, fileId string, nodes []string, exactSize int64) error {
	expectedShards, err := models.GetShardsByFileID(fileId)
	if err != nil {
		return fmt.Errorf("failed to retrieve shards: %v", err)
	}
	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return fmt.Errorf("failed to create Reed-Solomon encoder: %v", err)
	}
	shardFiles := make([]*os.File, TotalShards)
	for i := 0; i < TotalShards; i++ {
		path := filepath.Join(nodes[i], fmt.Sprintf("%s.shard.%d", fileId, i))
		if isValid(path, expectedShards[i].Checksum) {
			f, _ := os.Open(path)
			shardFiles[i] = f
			defer f.Close()
		} else {
			fmt.Printf("Shard %d is missing or CORRUPT. Marking for reconstruction.\n", i)
			shardFiles[i] = nil
		}
	}
	shardChunkSize := int64(ChunkSize)
	remaining := exactSize

	for remaining > 0 {
		shards := make([][]byte, TotalShards)
		for i := 0; i < TotalShards; i++ {
			if shardFiles[i] != nil {
				shards[i] = make([]byte, shardChunkSize)
				io.ReadFull(shardFiles[i], shards[i])
			}
		}
		if err = enc.Reconstruct(shards); err != nil {
			return fmt.Errorf("failed to reconstruct shards: %v", err)
		}

		toWrite := int64(ChunkSize / DataShards)
		if remaining < toWrite {
			toWrite = remaining
		}

		enc.Join(w, shards, int(toWrite))
		remaining -= toWrite
	}

	return nil
}

func isValid(path string, expectedHash string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false
	}

	actualHash := fmt.Sprintf("%x", h.Sum(nil))
	return actualHash == expectedHash
}
