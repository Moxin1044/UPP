package model

type Setting struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"column:setting_key;uniqueIndex;size:100;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}
