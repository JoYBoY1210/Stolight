package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/joyboy1210/stolight/models"
)

const ChunkSize = 10 * 1024 * 1024

func SplitFile(fileName string, fileSize int64, src io.Reader, nodes []string) error {
	fileId := uuid.New().String()
	fileRecord := models.File{
		ID:   fileId,
		Name: fileName,
		Size: fileSize,
	}
	err := models.CreateFile(&fileRecord)
	if err != nil {
		return fmt.Errorf("failed to create file record: %w", err)
	}
	buffer := make([]byte, ChunkSize)
	seq := 0
	for {
		bytesRead, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if bytesRead == 0 {
			break
		}
		nodeIdx := seq % len(nodes)
		targetNode := nodes[nodeIdx]
		chunkId := uuid.New().String()
		chunkFileName := fmt.Sprintf("chunk_%s_%d", chunkId, seq)
		chunkPath := filepath.Join(targetNode, chunkFileName)
		err = os.WriteFile(chunkPath, buffer[:bytesRead], 0644)
		if err != nil {
			return fmt.Errorf("failed to save to %s: %w", targetNode, err)
		}
		chunkRecord := models.ChunkMetaData{
			ID:            chunkId,
			FileID:        fileId,
			ChunkIndex:    seq,
			StorageNodeId: targetNode,
		}

		err = models.CreateChunkMetaData(&chunkRecord)
		if err != nil {
			return fmt.Errorf("failed to create chunk metadata: %w", err)
		}
		fmt.Printf("Chunk %d went to node %s\n", seq, targetNode)
		seq++
	}
	fmt.Println("file processed")
	return nil

}

func MergeFile(fileId string, dest io.Writer) error {
	var chunks []models.ChunkMetaData
	chunks, err := models.GetChunksByFileID(fileId)
	if err != nil {
		return fmt.Errorf("could not get chunks for file %s: %w", fileId, err)
	}
	for _, chunk := range chunks {
		chunkFileName := fmt.Sprintf("chunk_%s_%d", chunk.ID, chunk.ChunkIndex)
		chunkPath := filepath.Join(chunk.StorageNodeId, chunkFileName)

		file, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("could not open chunk file %s: %w", chunkPath, err)
		}
		defer file.Close()
		_, err = io.Copy(dest, file)

		if err != nil {
			return fmt.Errorf("could not write chunk data to dest: %w", err)
		}
	}
	return nil
}
