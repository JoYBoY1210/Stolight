package models

import "time"

type File struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"index"`
	Size      int64
	BucketID  string          `gorm:"index"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	Chunks    []ChunkMetaData `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE;"`
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
