package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/joyboy1210/stolight/models"
	"github.com/klauspost/reedsolomon"
	"github.com/zeebo/xxh3"
)

const (
	DataShards   = 4
	ParityShards = 2
	TotalShards  = DataShards + ParityShards
	ChunkSize    = 4 * 1024 * 1024
)

func EncodeFile(reader io.Reader, fileID string, NodeDirs []string) error {
	if len(NodeDirs) != TotalShards {
		return fmt.Errorf("expected %d node directories, got %d", TotalShards, len(NodeDirs))
	}

	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %s", err)
	}

	outFiles := make([]*os.File, TotalShards)
	shardIds := make([]string, TotalShards)
	for i := 0; i < TotalShards; i++ {
		shardIds[i] = uuid.New().String()
	}

	for i := 0; i < TotalShards; i++ {
		if err := os.MkdirAll(NodeDirs[i], 0755); err != nil {
			return cleanupFailedUpload(outFiles, fileID, NodeDirs, err)
		}

		tmpPath := filepath.Join(NodeDirs[i], fmt.Sprintf("%s.shard.%d.tmp", fileID, i))
		f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return cleanupFailedUpload(outFiles, fileID, NodeDirs, fmt.Errorf("failed to open shard file: %w", err))
		}
		outFiles[i] = f
	}

	var allSplits []models.Split
	var splitsMu sync.Mutex
	chunkIndex := 0

	buf := make([]byte, ChunkSize)
	for {
		n, err := io.ReadFull(reader, buf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return cleanupFailedUpload(outFiles, fileID, NodeDirs, fmt.Errorf("error reading stream: %w", err))
		}

		chunkData := buf[:n]
		shards, err := enc.Split(chunkData)
		if err != nil {
			return cleanupFailedUpload(outFiles, fileID, NodeDirs, fmt.Errorf("failed to split: %w", err))
		}
		if err := enc.Encode(shards); err != nil {
			return cleanupFailedUpload(outFiles, fileID, NodeDirs, fmt.Errorf("failed to encode: %w", err))
		}

		var wg sync.WaitGroup
		errCh := make(chan error, TotalShards)

		for i, shard := range shards {
			wg.Add(1)
			go func(i int, shard []byte, currentChunkIndex int) {
				defer wg.Done()

				if _, err := outFiles[i].Write(shard); err != nil {
					errCh <- fmt.Errorf("failed to write shard %d: %w", i, err)
					return
				}

				hashUint := xxh3.Hash(shard)
				checkSum := fmt.Sprintf("%x", hashUint)

				shardRecord := models.Split{
					ID:         uuid.New().String(),
					FileID:     fileID,
					ShardID:    shardIds[i],
					ChunkIndex: currentChunkIndex,
					Hash:       checkSum,
					Size:       int64(len(shard)),
				}

				splitsMu.Lock()
				allSplits = append(allSplits, shardRecord)
				splitsMu.Unlock()

			}(i, shard, chunkIndex)
		}

		wg.Wait()
		close(errCh)
		for err := range errCh {
			if err != nil {
				return cleanupFailedUpload(outFiles, fileID, NodeDirs, err)
			}
		}
		chunkIndex++
	}

	for _, f := range outFiles {
		if f != nil {
			f.Close()
		}
	}

	var shardRecords []models.Shard
	for i := 0; i < TotalShards; i++ {
		tmpPath := filepath.Join(NodeDirs[i], fmt.Sprintf("%s.shard.%d.tmp", fileID, i))
		finalPath := filepath.Join(NodeDirs[i], fmt.Sprintf("%s.shard.%d", fileID, i))

		if err := os.Rename(tmpPath, finalPath); err != nil {
			return err
		}

		shardRecords = append(shardRecords, models.Shard{
			Id:         shardIds[i],
			FileID:     fileID,
			ShardIndex: i,
			Path:       finalPath,
		})
	}

	if err := models.CreateShards(shardRecords); err != nil {
		return fmt.Errorf("failed to save shard metadata: %v", err)
	}

	if err := models.CreateBulkSplit(allSplits); err != nil {
		return fmt.Errorf("failed to save split metadata: %v", err)
	}

	return nil
}

func cleanupFailedUpload(openFiles []*os.File, fileID string, nodeDirs []string, originalErr error) error {
	for _, f := range openFiles {
		if f != nil {
			f.Close()
		}
	}
	for i, nodeDir := range nodeDirs {

		tmpPath := filepath.Join(nodeDir, fmt.Sprintf("%s.shard.%d.tmp", fileID, i))
		os.Remove(tmpPath)
	}
	return originalErr
}
