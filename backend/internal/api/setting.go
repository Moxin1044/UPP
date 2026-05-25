package api

import (
	"upp/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingService *service.SettingService
}

func NewSettingHandler(settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, err := h.settingService.GetAll()
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, settings)
}

type SetSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

func (h *SettingHandler) Set(c *gin.Context) {
	var req SetSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	if err := h.settingService.Set(req.Key, req.Value); err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, nil)
}
