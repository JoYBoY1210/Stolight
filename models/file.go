package models

import (
	"time"
)

type FileStatus string

const (
	FileStatusStaging   FileStatus = "staging"
	FileStatusPending   FileStatus = "pending"
	FileStatusEncoding  FileStatus = "encoding"
	FileStatusCompleted FileStatus = "completed"
	FileStatusFailed    FileStatus = "failed"
)

type File struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex:idx_file_name_bucket_id"`
	Size      int64
	BucketID  string     `gorm:"uniqueIndex:idx_file_name_bucket_id"`
	Status    FileStatus `gorm:"default:'staging';index:idx_file_status" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
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

// func DeleteFileByID(fileID string) error {
// 	err := DeleteChunksByFileID(fileID)
// 	if err != nil {
// 		return err
// 	}
// 	result := db.Delete(File{}, "id=?", fileID)
// 	if result.Error != nil {
// 		return fmt.Errorf("failed to delete the file and the chunks from db: %w", result.Error)
// 	}
// 	return nil
// }

func GetFileByFileNameAndBucketId(fileName string, bucketID string) (*File, error) {
	var file File
	result := db.Where("name = ? AND bucket_id = ?", fileName, bucketID).First(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}

func UpdateFileStatusAndSize(fileID string, status FileStatus, size int64) error {
	result := db.Model(&File{}).Where("id = ?", fileID).Updates(map[string]interface{}{"status": status, "size": size})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateFileStatus(fileID string, status FileStatus) error {
	result := db.Model(&File{}).Where("id = ?", fileID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
