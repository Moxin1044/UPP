package model

import "time"

type MonitorResult struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MonitorID uint      `gorm:"index;not null" json:"monitor_id"`
	Status    string    `gorm:"size:20;not null" json:"status"` // up, down
	Latency   float64   `json:"latency"`                        // ms
	Message   string    `gorm:"size:500" json:"message"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
