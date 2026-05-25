package api

import (
	"strconv"

	"upp/internal/model"
	"upp/internal/repository"
	"upp/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	monitorService    *service.MonitorService
	monitorResultRepo *repository.MonitorResultRepository
}

func NewDashboardHandler(monitorService *service.MonitorService, resultRepo *repository.MonitorResultRepository) *DashboardHandler {
	return &DashboardHandler{monitorService: monitorService, monitorResultRepo: resultRepo}
}

type DashboardStats struct {
	TotalMonitors int64   `json:"total_monitors"`
	UpMonitors    int64   `json:"up_monitors"`
	DownMonitors  int64   `json:"down_monitors"`
	AvgLatency    float64 `json:"avg_latency"`
}

type ProjectDashboard struct {
	Project       model.Project `json:"project"`
	TotalMonitors int           `json:"total_monitors"`
	Uptime        float64       `json:"uptime"`
	AvgLatency    float64       `json:"avg_latency"`
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	monitors, err := h.monitorService.ListByUser(userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	stats := DashboardStats{
		TotalMonitors: int64(len(monitors)),
	}

	var totalLatency float64
	var latencyCount float64

	for _, m := range monitors {
		latest, err := h.monitorResultRepo.FindLatestByMonitorID(m.ID)
		if err != nil {
			continue
		}
		if latest.Status == "up" {
			stats.UpMonitors++
		} else {
			stats.DownMonitors++
		}
		totalLatency += latest.Latency
		latencyCount++
	}

	if latencyCount > 0 {
		stats.AvgLatency = totalLatency / latencyCount
	}

	Success(c, stats)
}

func (h *DashboardHandler) GetProjectDashboard(c *gin.Context) {
	userID := c.GetUint("user_id")
	projectID, err := parseUint(c.Param("id"))
	if err != nil {
		BadRequest(c, "invalid project id")
		return
	}

	monitors, err := h.monitorService.ListByProject(uint(projectID))
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	var upCount, totalChecks int64
	var totalLatency float64
	var latencyCount float64

	for _, m := range monitors {
		if m.UserID != userID {
			continue
		}
		upC, _ := h.monitorResultRepo.CountByMonitorIDAndStatus(m.ID, "up")
		downC, _ := h.monitorResultRepo.CountByMonitorIDAndStatus(m.ID, "down")
		upCount += upC
		totalChecks += upC + downC

		avg, _ := h.monitorResultRepo.AvgLatencyByMonitorID(m.ID)
		if avg > 0 {
			totalLatency += avg
			latencyCount++
		}
	}

	var uptime float64
	if totalChecks > 0 {
		uptime = float64(upCount) / float64(totalChecks) * 100
	}
	var avgLatency float64
	if latencyCount > 0 {
		avgLatency = totalLatency / latencyCount
	}

	Success(c, gin.H{
		"total_monitors": len(monitors),
		"uptime":         uptime,
		"avg_latency":    avgLatency,
	})
}

func parseUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
