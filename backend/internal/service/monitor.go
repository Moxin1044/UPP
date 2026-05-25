package service

import (
	"upp/internal/model"
	"upp/internal/repository"
)

// SchedulerNotifier defines the interface for notifying the scheduler about monitor changes
type SchedulerNotifier interface {
	AddMonitor(m model.Monitor)
	UpdateMonitor(m model.Monitor)
	RemoveMonitor(monitorID uint)
}

type MonitorService struct {
	repo      *repository.MonitorRepository
	scheduler SchedulerNotifier
}

func NewMonitorService(repo *repository.MonitorRepository, scheduler SchedulerNotifier) *MonitorService {
	return &MonitorService{repo: repo, scheduler: scheduler}
}

type CreateMonitorRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	Name      string `json:"name" binding:"required,max=100"`
	Type      string `json:"type" binding:"required,oneof=http https tcp ping"`
	Target    string `json:"target" binding:"required,max=500"`
	Interval  int    `json:"interval"`
	Timeout   int    `json:"timeout"`
}

type UpdateMonitorRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Type     string `json:"type" binding:"required,oneof=http https tcp ping"`
	Target   string `json:"target" binding:"required,max=500"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
	Enabled  *bool  `json:"enabled"`
}

func (s *MonitorService) Create(userID uint, req CreateMonitorRequest) (*model.Monitor, error) {
	monitor := &model.Monitor{
		UserID:    userID,
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Type:      req.Type,
		Target:    req.Target,
		Interval:  req.Interval,
		Timeout:   req.Timeout,
		Enabled:   true,
	}
	if monitor.Interval <= 0 {
		monitor.Interval = 60
	}
	if monitor.Timeout <= 0 {
		monitor.Timeout = 10
	}
	if err := s.repo.Create(monitor); err != nil {
		return nil, err
	}
	if s.scheduler != nil && monitor.Enabled {
		s.scheduler.AddMonitor(*monitor)
	}
	return monitor, nil
}

func (s *MonitorService) ListByUser(userID uint) ([]model.Monitor, error) {
	return s.repo.FindByUserID(userID)
}

func (s *MonitorService) ListByProject(projectID uint) ([]model.Monitor, error) {
	return s.repo.FindByProjectID(projectID)
}

func (s *MonitorService) GetByID(id uint) (*model.Monitor, error) {
	return s.repo.FindByID(id)
}

func (s *MonitorService) Update(id uint, req UpdateMonitorRequest) (*model.Monitor, error) {
	monitor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	monitor.Name = req.Name
	monitor.Type = req.Type
	monitor.Target = req.Target
	monitor.Interval = req.Interval
	monitor.Timeout = req.Timeout
	if req.Enabled != nil {
		monitor.Enabled = *req.Enabled
	}
	if monitor.Interval <= 0 {
		monitor.Interval = 60
	}
	if monitor.Timeout <= 0 {
		monitor.Timeout = 10
	}
	if err := s.repo.Update(monitor); err != nil {
		return nil, err
	}
	if s.scheduler != nil {
		s.scheduler.UpdateMonitor(*monitor)
	}
	return monitor, nil
}

func (s *MonitorService) ToggleEnabled(id uint, enabled bool) (*model.Monitor, error) {
	monitor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	monitor.Enabled = enabled
	if err := s.repo.Update(monitor); err != nil {
		return nil, err
	}
	if s.scheduler != nil {
		s.scheduler.UpdateMonitor(*monitor)
	}
	return monitor, nil
}

func (s *MonitorService) Delete(id uint) error {
	if s.scheduler != nil {
		s.scheduler.RemoveMonitor(id)
	}
	return s.repo.Delete(id)
}

func (s *MonitorService) ListEnabled() ([]model.Monitor, error) {
	return s.repo.FindEnabled()
}
