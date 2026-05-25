package api

import (
	"strconv"

	"upp/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userService.ListAll()
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	user, err := h.userService.Create(req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, user)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	var req service.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	if err := h.userService.UpdatePassword(uint(id), req); err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	if err := h.userService.Delete(uint(id)); err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, nil)
}
