package service

import (
	"upp/internal/model"
	"upp/internal/repository"
)

type MonitorResultService struct {
	repo *repository.MonitorResultRepository
}

func NewMonitorResultService(repo *repository.MonitorResultRepository) *MonitorResultService {
	return &MonitorResultService{repo: repo}
}

func (s *MonitorResultService) Create(result *model.MonitorResult) error {
	return s.repo.Create(result)
}

func (s *MonitorResultService) ListByMonitor(monitorID uint, limit int) ([]model.MonitorResult, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.FindByMonitorID(monitorID, limit)
}

func (s *MonitorResultService) GetLatest(monitorID uint) (*model.MonitorResult, error) {
	return s.repo.FindLatestByMonitorID(monitorID)
}

type MonitorStats struct {
	TotalChecks int64   `json:"total_checks"`
	UpChecks    int64   `json:"up_checks"`
	DownChecks  int64   `json:"down_checks"`
	Uptime      float64 `json:"uptime"` // percentage
	AvgLatency  float64 `json:"avg_latency"`
}

func (s *MonitorResultService) GetStats(monitorID uint) (*MonitorStats, error) {
	upCount, err := s.repo.CountByMonitorIDAndStatus(monitorID, "up")
	if err != nil {
		return nil, err
	}
	downCount, err := s.repo.CountByMonitorIDAndStatus(monitorID, "down")
	if err != nil {
		return nil, err
	}
	avgLatency, err := s.repo.AvgLatencyByMonitorID(monitorID)
	if err != nil {
		return nil, err
	}
	total := upCount + downCount
	var uptime float64
	if total > 0 {
		uptime = float64(upCount) / float64(total) * 100
	}
	return &MonitorStats{
		TotalChecks: total,
		UpChecks:    upCount,
		DownChecks:  downCount,
		Uptime:      uptime,
		AvgLatency:  avgLatency,
	}, nil
}
