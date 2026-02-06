package models

import "time"

type User struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	UserName       string    `gorm:"uniqueIndex" json:"username"`
	PasswordHash   string    `json:"-"`
	Key            string    `gorm:"uniqueIndex" json:"key"`
	Role           string    `gorm:"index" json:"role"`
	AllowedBuckets string    `json:"allowed_buckets"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}
