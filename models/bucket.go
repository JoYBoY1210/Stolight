package models

import (
	"time"

	"github.com/google/uuid"
)

type Bucket struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Files     []File    `gorm:"foreignKey:BucketID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
}

func CreateBucket(name string) (*Bucket, error) {
	bucket := Bucket{
		ID:   uuid.NewString(),
		Name: name,
	}
	result := db.Create(&bucket)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bucket, nil
}

func GetBucketByName(name string) (*Bucket, error) {
	var bucket Bucket
	result := db.First(&bucket, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bucket, nil
}
