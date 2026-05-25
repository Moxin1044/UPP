package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"upp/internal/model"
	"upp/internal/probe"
	"upp/internal/repository"
)

type Scheduler struct {
	monitorRepo      *repository.MonitorRepository
	resultRepo       *repository.MonitorResultRepository
	running          map[uint]context.CancelFunc
	mu               sync.Mutex
	notificationChan chan NotificationEvent
}

type NotificationEvent struct {
	Monitor   model.Monitor
	OldStatus string
	NewStatus string
	Result    model.MonitorResult
}

func NewScheduler(monitorRepo *repository.MonitorRepository, resultRepo *repository.MonitorResultRepository, notifyChan chan NotificationEvent) *Scheduler {
	return &Scheduler{
		monitorRepo:      monitorRepo,
		resultRepo:       resultRepo,
		running:          make(map[uint]context.CancelFunc),
		notificationChan: notifyChan,
	}
}

func (s *Scheduler) Start() {
	log.Println("[Scheduler] Starting...")
	s.loadMonitors()
}

func (s *Scheduler) loadMonitors() {
	monitors, err := s.monitorRepo.FindEnabled()
	if err != nil {
		log.Printf("[Scheduler] Error loading monitors: %v", err)
		return
	}
	for _, m := range monitors {
		s.AddMonitor(m)
	}
}

func (s *Scheduler) AddMonitor(m model.Monitor) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.running[m.ID]; exists {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.running[m.ID] = cancel

	go s.runMonitor(ctx, m)
	log.Printf("[Scheduler] Added monitor: %s (%s)", m.Name, m.Type)
}

func (s *Scheduler) RemoveMonitor(monitorID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cancel, exists := s.running[monitorID]; exists {
		cancel()
		delete(s.running, monitorID)
		log.Printf("[Scheduler] Removed monitor: %d", monitorID)
	}
}

func (s *Scheduler) UpdateMonitor(m model.Monitor) {
	s.RemoveMonitor(m.ID)
	if m.Enabled {
		s.AddMonitor(m)
	}
}

func (s *Scheduler) runMonitor(ctx context.Context, m model.Monitor) {
	interval := time.Duration(m.Interval) * time.Second
	if interval < time.Second {
		interval = time.Second
	}

	var runningMu sync.Mutex
	isRunning := false

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately on start
	s.probeMonitor(&m, &runningMu, &isRunning)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.probeMonitor(&m, &runningMu, &isRunning)
		}
	}
}

func (s *Scheduler) probeMonitor(m *model.Monitor, mu *sync.Mutex, isRunning *bool) {
	mu.Lock()
	if *isRunning {
		mu.Unlock()
		return // Skip if previous probe still running
	}
	*isRunning = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		*isRunning = false
		mu.Unlock()
	}()

	p, ok := probe.Get(m.Type)
	if !ok {
		log.Printf("[Scheduler] Unknown probe type: %s", m.Type)
		return
	}

	timeout := time.Duration(m.Timeout) * time.Second
	if timeout < time.Second {
		timeout = 10 * time.Second
	}

	result := p.Probe(m.Target, timeout)

	// Get old status for notification
	var oldStatus string
	latest, err := s.resultRepo.FindLatestByMonitorID(m.ID)
	if err == nil {
		oldStatus = latest.Status
	}

	monitorResult := &model.MonitorResult{
		MonitorID: m.ID,
		Status:    result.Status,
		Latency:   result.Latency,
		Message:   result.Message,
		CreatedAt: time.Now(),
	}

	if err := s.resultRepo.Create(monitorResult); err != nil {
		log.Printf("[Scheduler] Error saving result for monitor %d: %v", m.ID, err)
		return
	}

	// Send notification if status changed
	if oldStatus != "" && oldStatus != result.Status && s.notificationChan != nil {
		select {
		case s.notificationChan <- NotificationEvent{
			Monitor:   *m,
			OldStatus: oldStatus,
			NewStatus: result.Status,
			Result:    *monitorResult,
		}:
		default:
		}
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, cancel := range s.running {
		cancel()
		delete(s.running, id)
	}
	log.Println("[Scheduler] Stopped")
}
