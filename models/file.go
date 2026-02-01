package models

import "time"

type File struct {
	ID         string `gorm:"primaryKey"`
	Name       string `gorm:"index"`
	Size       int64
	Created_at time.Time       `gorm:"autoCreateTime"`
	Updated_at time.Time       `gorm:"autoUpdateTime"`
	Chunks     []ChunkMetaData `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE;"`
}

type ChunkMetaData struct {
	ID            string `gorm:"primaryKey"`
	FileID        string `gorm:"index"` 
	ChunkIndex    int
	StorageNodeId string    `gorm:"index"`
	Created_at    time.Time `gorm:"autoCreateTime"`
	Updated_at    time.Time `gorm:"autoUpdateTime"`
}
