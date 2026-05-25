package api

import (
	"strconv"

	"upp/internal/service"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	monitorService      *service.MonitorService
	monitorResultService *service.MonitorResultService
}

func NewMonitorHandler(monitorService *service.MonitorService, resultService *service.MonitorResultService) *MonitorHandler {
	return &MonitorHandler{monitorService: monitorService, monitorResultService: resultService}
}

func (h *MonitorHandler) Create(c *gin.Context) {
	var req service.CreateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	monitor, err := h.monitorService.Create(userID, req)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, monitor)
}

func (h *MonitorHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	monitors, err := h.monitorService.ListByUser(userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, monitors)
}

func (h *MonitorHandler) ListByProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid project id")
		return
	}
	monitors, err := h.monitorService.ListByProject(uint(projectID))
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, monitors)
}

func (h *MonitorHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	monitor, err := h.monitorService.GetByID(uint(id))
	if err != nil {
		NotFound(c, "monitor not found")
		return
	}
	Success(c, monitor)
}

func (h *MonitorHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	var req service.UpdateMonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	monitor, err := h.monitorService.Update(uint(id), req)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, monitor)
}

func (h *MonitorHandler) ToggleEnabled(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		BadRequest(c, err.Error())
		return
	}
	monitor, err := h.monitorService.ToggleEnabled(uint(id), body.Enabled)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, monitor)
}

func (h *MonitorHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	if err := h.monitorService.Delete(uint(id)); err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, nil)
}

func (h *MonitorHandler) GetResults(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	results, err := h.monitorResultService.ListByMonitor(uint(id), limit)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, results)
}

func (h *MonitorHandler) GetStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	stats, err := h.monitorResultService.GetStats(uint(id))
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, stats)
}
