package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joyboy1210/stolight/models"
	"github.com/klauspost/reedsolomon"
	"github.com/zeebo/xxh3"
)

func DecodeFile(w io.Writer, fileId string, nodes []string, exactSize int64) error {
	expectedShards, err := models.GetShardsByFileID(fileId)
	if err != nil {
		return fmt.Errorf("failed to retrieve shards: %v", err)
	}

	allSplits, err := models.GetSplitsByFileID(fileId)
	if err != nil {
		return fmt.Errorf("failed to retrieve splits: %v", err)
	}

	shardIndexMap := make(map[string]int)
	for _, s := range expectedShards {
		shardIndexMap[s.Id] = s.ShardIndex
	}

	type splitData struct {
		Hash string
		Size int64
	}

	chunkHashes := make(map[int]map[int]splitData)

	for _, sp := range allSplits {
		sIdx, exists := shardIndexMap[sp.ShardID]
		if !exists {
			continue
		}

		if chunkHashes[sp.ChunkIndex] == nil {
			chunkHashes[sp.ChunkIndex] = make(map[int]splitData)
		}

		chunkHashes[sp.ChunkIndex][sIdx] = splitData{
			Hash: sp.Hash,
			Size: sp.Size,
		}
	}

	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return fmt.Errorf("failed to create Reed-Solomon encoder: %v", err)
	}

	shardFiles := make([]*os.File, TotalShards)

	for i := 0; i < TotalShards; i++ {
		path := filepath.Join(nodes[i], fmt.Sprintf("%s.shard.%d", fileId, i))

		f, err := os.Open(path)
		if err == nil {
			shardFiles[i] = f
			defer f.Close()
		}
	}

	maxShardSize := int64(ChunkSize / DataShards)

	buffers := make([][]byte, TotalShards)
	for i := range buffers {
		buffers[i] = make([]byte, maxShardSize)
	}

	remaining := exactSize
	chunkIndex := 0

	for remaining > 0 {
		stepShards := make([][]byte, TotalShards)

		for i := 0; i < TotalShards; i++ {
			if shardFiles[i] == nil {
				continue
			}

			expected := chunkHashes[chunkIndex][i]

			if expected.Size == 0 {
				shardFiles[i] = nil
				continue
			}

			stepShards[i] = buffers[i][:expected.Size]

			_, err := io.ReadFull(shardFiles[i], stepShards[i])
			if err != nil {
				shardFiles[i] = nil
				stepShards[i] = nil
				continue
			}

			actualHash := fmt.Sprintf("%x", xxh3.Hash(stepShards[i]))

			if actualHash != expected.Hash {
				shardFiles[i] = nil
				stepShards[i] = nil
			}
		}

		if err := enc.Reconstruct(stepShards); err != nil {
			return fmt.Errorf("failed to reconstruct chunk %d: %v", chunkIndex, err)
		}

		outSize := int64(ChunkSize)
		if remaining < outSize {
			outSize = remaining
		}

		if err := enc.Join(w, stepShards, int(outSize)); err != nil {
			return fmt.Errorf("failed to join chunk %d: %v", chunkIndex, err)
		}

		remaining -= outSize
		chunkIndex++
	}

	return nil
}
