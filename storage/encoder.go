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
)

const (
	DataShards   = 4
	ParityShards = 2
	TotalShards  = DataShards + ParityShards
	ChunkSize    = 4 * 1024 * 1024
)



func EncodeFile(reader io.Reader, fileName string, NodeDirs []string, fileSize int64, bucketName string) error {
	if len(NodeDirs) != TotalShards {
		return fmt.Errorf("expected %d node directories, got %d", TotalShards, len(NodeDirs))
	}
	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %s", err)
	}
	outFiles := make([]*os.File, TotalShards)
	storageName := fmt.Sprintf("%s_%s", bucketName, fileName)
	for i := 0; i < TotalShards; i++ {
		if err := os.MkdirAll(NodeDirs[i], 0755); err != nil {
			return cleanupFailedUpload(outFiles, storageName, NodeDirs, fmt.Errorf("failed to create node dir  %s", err))
		}
		outPath := filepath.Join(NodeDirs[i], fmt.Sprintf("%s.shard.%d", storageName, i))
		f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return cleanupFailedUpload(outFiles, storageName, NodeDirs, fmt.Errorf("failed to open shard file: %w", err))
		}
		outFiles[i] = f

	}
	buf := make([]byte, ChunkSize)
	var totalSize int64
	for {
		n, err := io.ReadFull(reader, buf)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return cleanupFailedUpload(outFiles, storageName, NodeDirs, fmt.Errorf("error reading stream: %w", err))
		}
		totalSize += int64(n)
		chunkData := buf[:n]
		shards, err := enc.Split(chunkData)
		if err != nil {
			return cleanupFailedUpload(outFiles, storageName, NodeDirs, fmt.Errorf("failed to split data into shards: %w", err))
		}
		if err := enc.Encode(shards); err != nil {
			return cleanupFailedUpload(outFiles, storageName, NodeDirs, fmt.Errorf("failed to encode shards: %w", err))
		}
		var wg sync.WaitGroup
		errCh := make(chan error, TotalShards)
		for i, shard := range shards {
			wg.Add(1)
			go func(i int, shard []byte) {
				defer wg.Done()
				if _, err := outFiles[i].Write(shard); err != nil {
					errCh <- fmt.Errorf("failed to write shard %d: %w", i, err)
				}
			}(i, shard)
		}
		wg.Wait()
		close(errCh)
		for err := range errCh {
			if err != nil {
				return cleanupFailedUpload(outFiles, storageName, NodeDirs, err)
			}
		}
	}
	for _, f := range outFiles {
		if f != nil {
			f.Close()
		}
	}

	bucketId, err := models.GetBucketByName(bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket id: %w", err)
	}

	fileId := uuid.New().String()
	fileRecord := models.File{
		ID:       fileId,
		Name:     fileName,
		Size:     fileSize,
		BucketID: bucketId.ID,
	}
	err = models.CreateFile(&fileRecord)
	if err != nil {
		return fmt.Errorf("failed to create file record: %w", err)
	}

	return nil

}

func cleanupFailedUpload(openFiles []*os.File, storageName string, nodeDirs []string, originalErr error) error {

	for _, f := range openFiles {
		if f != nil {
			f.Close()
		}
	}
	for i, nodeDir := range nodeDirs {
		outPath := filepath.Join(nodeDir, fmt.Sprintf("%s.shard.%d", storageName, i))
		os.Remove(outPath)
	}

	return originalErr
}
