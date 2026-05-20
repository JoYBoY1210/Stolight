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
			f, err := os.Open(path)
			if err == nil {
				shardFiles[i] = f
				defer f.Close()
			} else {
				shardFiles[i] = nil
			}
		} else {
			fmt.Printf("Shard %d is missing or CORRUPT. Marking for reconstruction.\n", i)
			shardFiles[i] = nil
		}
	}

	shardChunkSize := int64(ChunkSize / DataShards)
	remaining := exactSize

	buffers := make([][]byte, TotalShards)
	for i := range buffers {
		buffers[i] = make([]byte, shardChunkSize)
	}

	for remaining > 0 {
		stepShards := make([][]byte, TotalShards)

		for i := 0; i < TotalShards; i++ {
			if shardFiles[i] != nil {
				stepShards[i] = buffers[i]

				_, err := io.ReadFull(shardFiles[i], stepShards[i])
				if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
					fmt.Printf("Error reading shard %d midway: %v. Dropping shard.\n", i, err)
					shardFiles[i] = nil
					stepShards[i] = nil
				}
			} else {
				stepShards[i] = nil
			}
		}

		if err = enc.Reconstruct(stepShards); err != nil {
			return fmt.Errorf("failed to reconstruct shards: %v", err)
		}
		outSize := int64(ChunkSize)
		if remaining < outSize {
			outSize = remaining
		}

		if err := enc.Join(w, stepShards, int(outSize)); err != nil {
			return fmt.Errorf("failed to join shards: %v", err)
		}

		remaining -= outSize
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
