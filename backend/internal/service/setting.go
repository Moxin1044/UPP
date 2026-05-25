package service

import (
	"upp/internal/repository"
)

type SettingService struct {
	repo *repository.SettingRepository
}

func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) Get(key string) (string, error) {
	return s.repo.Get(key)
}

func (s *SettingService) Set(key, value string) error {
	return s.repo.Set(key, value)
}

func (s *SettingService) GetAll() (map[string]string, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}
