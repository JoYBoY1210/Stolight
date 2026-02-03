package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/joyboy1210/stolight/models"
	"gorm.io/gorm"
)

const ChunkSize = 10 * 1024 * 1024

func SplitFile(db *gorm.DB, fileName string, fileSize int64, src io.Reader, nodes []string) error {
	fileId := uuid.New().String()
	fileRecord := models.File{
		ID:   fileId,
		Name: fileName,
		Size: fileSize,
	}
	if err := db.Create(&fileRecord).Error; err != nil {
		return err
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

		if err := db.Create(&chunkRecord).Error; err != nil {
			return fmt.Errorf("failed to save chunk meta: %w", err)
		}
		fmt.Printf("Chunk %d went to node %s\n", seq, targetNode)
		seq++
	}
	fmt.Println("file processed")
	return nil

}

func MergeFile(db *gorm.DB, fileId string, dest io.Writer) error {
	var chunks []models.ChunkMetaData
	if err := db.Where("file_id=?", fileId).Order("chunk_index asc").Find(&chunks); err.Error != nil {

		return fmt.Errorf("could not get the file parts from DB")
	}
	for _, chunk := range chunks {
		chunkFileName := fmt.Sprintf("chunk_%s_%d", chunk.ID, chunk.ChunkIndex)
		chunkPath := filepath.Join(chunk.StorageNodeId, chunkFileName)

		data, err := os.ReadFile(chunkPath)
		if err != nil {
			return fmt.Errorf("could not read chunk file %s: %w", chunkPath, err)
		}
		_, err = dest.Write(data)
		if err != nil {
			return fmt.Errorf("could not write chunk data to dest: %w", err)
		}
	}
	return nil
}
