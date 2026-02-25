package models

type Shard struct {
	Id       string `gorm:"primaryKey"`
	FileID   string `gorm:"index"`
	Index    int
	Path     string
	Checksum string
}

func CreateShards(shards []Shard) error {
	result := db.Create(&shards)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
