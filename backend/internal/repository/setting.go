package repository

import (
	"upp/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) Get(key string) (string, error) {
	var setting model.Setting
	if err := r.db.Where("setting_key = ?", key).First(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (r *SettingRepository) Set(key, value string) error {
	setting := model.Setting{Key: key, Value: value}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&setting).Error
}

func (r *SettingRepository) GetAll() ([]model.Setting, error) {
	var settings []model.Setting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}
