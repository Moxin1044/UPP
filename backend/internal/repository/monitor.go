package repository

import (
	"upp/internal/model"

	"gorm.io/gorm"
)

type MonitorRepository struct {
	db *gorm.DB
}

func NewMonitorRepository(db *gorm.DB) *MonitorRepository {
	return &MonitorRepository{db: db}
}

func (r *MonitorRepository) Create(monitor *model.Monitor) error {
	return r.db.Create(monitor).Error
}

func (r *MonitorRepository) FindByUserID(userID uint) ([]model.Monitor, error) {
	var monitors []model.Monitor
	if err := r.db.Where("user_id = ?", userID).Find(&monitors).Error; err != nil {
		return nil, err
	}
	return monitors, nil
}

func (r *MonitorRepository) FindByProjectID(projectID uint) ([]model.Monitor, error) {
	var monitors []model.Monitor
	if err := r.db.Where("project_id = ?", projectID).Find(&monitors).Error; err != nil {
		return nil, err
	}
	return monitors, nil
}

func (r *MonitorRepository) FindByID(id uint) (*model.Monitor, error) {
	var monitor model.Monitor
	if err := r.db.First(&monitor, id).Error; err != nil {
		return nil, err
	}
	return &monitor, nil
}

func (r *MonitorRepository) FindEnabled() ([]model.Monitor, error) {
	var monitors []model.Monitor
	if err := r.db.Where("enabled = ?", true).Find(&monitors).Error; err != nil {
		return nil, err
	}
	return monitors, nil
}

func (r *MonitorRepository) Update(monitor *model.Monitor) error {
	return r.db.Save(monitor).Error
}

func (r *MonitorRepository) Delete(id uint) error {
	return r.db.Delete(&model.Monitor{}, id).Error
}
