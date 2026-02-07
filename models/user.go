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

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByKey(key string) (*User, error) {
	var user User
	err := db.Where("key = ?", key).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(Id, name, key, role, allowedBuckets string) error {
	user := User{
		ID:             Id,
		Username:       name,
		Key:            key,
		Role:           role,
		AllowedBuckets: allowedBuckets,
	}
	return db.Create(&user).Error
}
