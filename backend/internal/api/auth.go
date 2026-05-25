package api

import (
	"upp/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	resp, err := h.authService.Login(req)
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}
	Success(c, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	user, err := h.authService.Register(req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, user)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.authService.UserByID(userID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}
	Success(c, user)
}
