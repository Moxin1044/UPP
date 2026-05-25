package model

import "time"

type Monitor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ProjectID uint      `gorm:"index;not null" json:"project_id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Type      string    `gorm:"size:20;not null" json:"type"` // http, https, tcp, ping
	Target    string    `gorm:"size:500;not null" json:"target"`
	Interval  int       `gorm:"default:60" json:"interval"`  // seconds
	Timeout   int       `gorm:"default:10" json:"timeout"`   // seconds
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
