package models

import (
	"time"
)

type ChunkMetaData struct {
	ID            string `gorm:"primaryKey"`
	FileID        string `gorm:"index"`
	ChunkIndex    int
	StorageNodeId string    `gorm:"index"`
	Created_at    time.Time `gorm:"autoCreateTime"`
	Updated_at    time.Time `gorm:"autoUpdateTime"`
}

func CreateChunkMetaData(chunkData *ChunkMetaData) error {
	result := db.Create(&chunkData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetChunksByFileID(fileID string) ([]ChunkMetaData, error) {
	var chunks []ChunkMetaData
	result:= db.Where("file_id = ?", fileID).Order("chunk_index asc").Find(&chunks)
	if result.Error != nil {
		return nil, result.Error
	}
	return chunks, nil
}
