package models

type Split struct {
	ID         string `gorm:"primaryKey" json:"id"`
	FileID     string `gorm:"index" json:"file_id"`
	ShardID    string `json:"shard_id"`
	ChunkIndex int    `json:"index"`

	Hash string `json:"hash"`

	Size int64 `json:"size"`
}

func CreateBulkSplit(split []Split) error {
	return db.Create(split).Error
}

func GetSplitsByFileID(fileId string) ([]Split, error) {
	var splits []Split
	result := db.Where("file_id = ?", fileId).Find(&splits)
	if result.Error != nil {
		return nil, result.Error
	}
	return splits, nil
}
