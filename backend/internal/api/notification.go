package api

import (
	"strconv"

	"upp/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) Create(c *gin.Context) {
	var req service.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	notification, err := h.notificationService.Create(userID, req)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, notification)
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	notifications, err := h.notificationService.ListByUser(userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, notifications)
}

func (h *NotificationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	notification, err := h.notificationService.GetByID(uint(id))
	if err != nil {
		NotFound(c, "notification not found")
		return
	}
	Success(c, notification)
}

func (h *NotificationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	var req service.UpdateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	notification, err := h.notificationService.Update(uint(id), req)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, notification)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	if err := h.notificationService.Delete(uint(id)); err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, nil)
}

func (h *NotificationHandler) Test(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	notification, err := h.notificationService.GetByID(uint(id))
	if err != nil {
		NotFound(c, "notification not found")
		return
	}
	// Send test notification
	if err := h.notificationService.SendTest(notification); err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, nil)
}
