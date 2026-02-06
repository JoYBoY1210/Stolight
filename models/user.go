package models

import "time"

type User struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"uniqueIndex" json:"username"`
	PasswordHash   string    `json:"-"`
	Key            string    `gorm:"uniqueIndex" json:"key"`
	Role           string    `gorm:"index" json:"role"`
	AllowedBuckets string    `json:"allowed_buckets"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func GetTotalUsers() (int64, error) {
	var count int64
	err := db.Model(&User{}).Count(&count).Error
	return count, err
}
