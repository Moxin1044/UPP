package repository

import (
	"upp/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) FindByUserID(userID uint) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := r.db.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) FindByID(id uint) (*model.Notification, error) {
	var notification model.Notification
	if err := r.db.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) FindEnabledByUserID(userID uint) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := r.db.Where("user_id = ? AND enabled = ?", userID, true).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) Update(notification *model.Notification) error {
	return r.db.Save(notification).Error
}

func (r *NotificationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Notification{}, id).Error
}
