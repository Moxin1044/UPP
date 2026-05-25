package repository

import (
	"upp/internal/model"

	"gorm.io/gorm"
)

type MonitorResultRepository struct {
	db *gorm.DB
}

func NewMonitorResultRepository(db *gorm.DB) *MonitorResultRepository {
	return &MonitorResultRepository{db: db}
}

func (r *MonitorResultRepository) Create(result *model.MonitorResult) error {
	return r.db.Create(result).Error
}

func (r *MonitorResultRepository) FindByMonitorID(monitorID uint, limit int) ([]model.MonitorResult, error) {
	var results []model.MonitorResult
	if err := r.db.Where("monitor_id = ?", monitorID).Order("created_at DESC").Limit(limit).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MonitorResultRepository) FindLatestByMonitorID(monitorID uint) (*model.MonitorResult, error) {
	var result model.MonitorResult
	if err := r.db.Where("monitor_id = ?", monitorID).Order("created_at DESC").First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *MonitorResultRepository) CountByMonitorIDAndStatus(monitorID uint, status string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.MonitorResult{}).Where("monitor_id = ? AND status = ?", monitorID, status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MonitorResultRepository) AvgLatencyByMonitorID(monitorID uint) (float64, error) {
	var avg float64
	if err := r.db.Model(&model.MonitorResult{}).Where("monitor_id = ?", monitorID).Select("AVG(latency)").Scan(&avg).Error; err != nil {
		return 0, err
	}
	return avg, nil
}
