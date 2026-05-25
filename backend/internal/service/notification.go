package service

import (
	"encoding/json"
	"upp/internal/model"
	"upp/internal/repository"
)

type NotificationService struct {
	repo       *repository.NotificationRepository
	dispatcher NotificationDispatcher
}

type NotificationDispatcher interface {
	SendTest(n *model.Notification) error
}

func NewNotificationService(repo *repository.NotificationRepository, dispatcher NotificationDispatcher) *NotificationService {
	return &NotificationService{repo: repo, dispatcher: dispatcher}
}

type CreateNotificationRequest struct {
	Type    string                 `json:"type" binding:"required,oneof=lark dingtalk webhook email"`
	Config  map[string]interface{} `json:"config" binding:"required"`
	Enabled *bool                  `json:"enabled"`
}

type UpdateNotificationRequest struct {
	Type    string                 `json:"type" binding:"required,oneof=lark dingtalk webhook email"`
	Config  map[string]interface{} `json:"config" binding:"required"`
	Enabled *bool                  `json:"enabled"`
}

func (s *NotificationService) Create(userID uint, req CreateNotificationRequest) (*model.Notification, error) {
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return nil, err
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	notification := &model.Notification{
		UserID:  userID,
		Type:    req.Type,
		Config:  string(configBytes),
		Enabled: enabled,
	}
	if err := s.repo.Create(notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) ListByUser(userID uint) ([]model.Notification, error) {
	return s.repo.FindByUserID(userID)
}

func (s *NotificationService) GetByID(id uint) (*model.Notification, error) {
	return s.repo.FindByID(id)
}

func (s *NotificationService) Update(id uint, req UpdateNotificationRequest) (*model.Notification, error) {
	notification, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return nil, err
	}
	notification.Type = req.Type
	notification.Config = string(configBytes)
	if req.Enabled != nil {
		notification.Enabled = *req.Enabled
	}
	if err := s.repo.Update(notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *NotificationService) GetEnabledByUser(userID uint) ([]model.Notification, error) {
	return s.repo.FindEnabledByUserID(userID)
}

func (s *NotificationService) SendTest(n *model.Notification) error {
	return s.dispatcher.SendTest(n)
}
