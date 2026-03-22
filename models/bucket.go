package models

import (
	"fmt"
	"strings"
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
func GetBucketByID(id string) (*Bucket, error) {
	var bucket Bucket
	result := db.First(&bucket, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bucket, nil
}

func ValidateBuckets(raw string) (bool, error) {
	parts := strings.Split(raw, ",")

	unique := make(map[string]struct{})
	var names []string

	for _, p := range parts {
		name := strings.TrimSpace(p)
		if name == "" {
			continue
		}

		name = strings.ToLower(name)

		if _, exists := unique[name]; !exists {
			unique[name] = struct{}{}
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		return false, fmt.Errorf("no valid bucket names provided")
	}

	var count int64
	err := db.Model(&Bucket{}).
		Where("LOWER(name) IN ?", names).
		Distinct("name").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == int64(len(names)), nil
}
