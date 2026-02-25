package storage

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/joyboy1210/stolight/models"
	"github.com/klauspost/reedsolomon"
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
	hashers := make([]hash.Hash, TotalShards)
	writers := make([]io.Writer, TotalShards)

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
		hashers[i] = sha256.New()
		writers[i] = io.MultiWriter(f, hashers[i])
	}

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
			go func(i int, shard []byte) {
				defer wg.Done()

				if _, err := writers[i].Write(shard); err != nil {
					errCh <- fmt.Errorf("failed to write shard %d: %w", i, err)
				}
			}(i, shard)
		}
		wg.Wait()
		close(errCh)
		for err := range errCh {
			if err != nil {
				return cleanupFailedUpload(outFiles, fileID, NodeDirs, err)
			}
		}
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

		checksum := fmt.Sprintf("%x", hashers[i].Sum(nil))

		shardRecords = append(shardRecords, models.Shard{
			Id:       uuid.New().String(),
			FileID:   fileID,
			Index:    i,
			Path:     finalPath,
			Checksum: checksum,
		})
	}

	if err := models.CreateShards(shardRecords); err != nil {
		return fmt.Errorf("failed to save shard metadata: %v", err)
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
