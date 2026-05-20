package models

type Split struct {
	ID      string `gorm:"primaryKey" json:"id"`
	FileID  string `gorm:"index" json:"file_id"`
	ShardID string `json:"shard_id"`
	ChunkIndex   int    `json:"index"`

	Hash uint64 `json:"hash"`

	Size int64 `json:"size"`
}

func CreateBulkSplit(split []Split) error {
	return db.Create(split).Error
}
