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

func CreateFile(file *File) error {
	result := db.Create(&file)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetFileByID(fileID string) (*File, error) {
	var file File
	result := db.First(&file, "id = ?", fileID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}
