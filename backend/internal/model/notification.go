package model

type Notification struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	Type    string `gorm:"size:20;not null" json:"type"` // lark, dingtalk, webhook, email
	Config  string `gorm:"type:text;not null" json:"config"`
	Enabled bool   `gorm:"default:true" json:"enabled"`
}
