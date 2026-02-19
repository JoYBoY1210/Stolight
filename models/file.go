package models

import (
	"fmt"
	"time"
)

type File struct {
	ID           string `gorm:"primaryKey"`
	Name         string `gorm:"uniqueIndex:idx_file_name_bucket_id"`
	Size         int64
	BucketID     string          `gorm:"uniqueIndex:idx_file_name_bucket_id"`
	CreatedAt    time.Time       `gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime"`
	Chunks       []ChunkMetaData `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE;"`
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

func GetFilesByBucketID(bucketID string) ([]File, error) {
	var files []File
	result := db.Where("bucket_id = ?", bucketID).Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}
	return files, nil
}

func DeleteFileByID(fileID string) error {
	err := DeleteChunksByFileID(fileID)
	if err != nil {
		return err
	}
	result := db.Delete(File{}, "id=?", fileID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete the file and the chunks from db: %w", result.Error)
	}
	return nil
}

func GetFileByFileNameAndBucketId(fileName string, bucketID string) (*File, error) {
	var file File
	result := db.Where("name = ? AND bucket_id = ?", fileName, bucketID).First(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}
