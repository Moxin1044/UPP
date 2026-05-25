package api

import (
	"upp/internal/middleware"
	"upp/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *AuthHandler,
	projectHandler *ProjectHandler,
	monitorHandler *MonitorHandler,
	notificationHandler *NotificationHandler,
	userHandler *UserHandler,
	settingHandler *SettingHandler,
	dashboardHandler *DashboardHandler,
	hub *websocket.Hub,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)

		// Authenticated routes
		auth := api.Group("")
		auth.Use(middleware.JWTAuth(jwtSecret))
		{
			// Auth
			auth.GET("/auth/profile", authHandler.GetProfile)

			// Dashboard
			auth.GET("/dashboard/stats", dashboardHandler.GetStats)
			auth.GET("/dashboard/project/:id", dashboardHandler.GetProjectDashboard)

			// Projects
			auth.GET("/projects", projectHandler.List)
			auth.POST("/projects", projectHandler.Create)
			auth.GET("/projects/:id", projectHandler.Get)
			auth.PUT("/projects/:id", projectHandler.Update)
			auth.DELETE("/projects/:id", projectHandler.Delete)

			// Monitors
			auth.GET("/monitors", monitorHandler.List)
			auth.POST("/monitors", monitorHandler.Create)
			auth.GET("/monitors/:id", monitorHandler.Get)
			auth.PUT("/monitors/:id", monitorHandler.Update)
			auth.PATCH("/monitors/:id/toggle", monitorHandler.ToggleEnabled)
			auth.DELETE("/monitors/:id", monitorHandler.Delete)
			auth.GET("/monitors/:id/results", monitorHandler.GetResults)
			auth.GET("/monitors/:id/stats", monitorHandler.GetStats)
			auth.GET("/projects/:id/monitors", monitorHandler.ListByProject)

			// Notifications
			auth.GET("/notifications", notificationHandler.List)
			auth.POST("/notifications", notificationHandler.Create)
			auth.GET("/notifications/:id", notificationHandler.Get)
			auth.PUT("/notifications/:id", notificationHandler.Update)
			auth.DELETE("/notifications/:id", notificationHandler.Delete)
			auth.POST("/notifications/:id/test", notificationHandler.Test)

			// WebSocket
			auth.GET("/ws", hub.HandleConnection)

			// Admin routes
			admin := auth.Group("")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/users", userHandler.List)
				admin.POST("/users", userHandler.Create)
				admin.PUT("/users/:id/password", userHandler.UpdatePassword)
				admin.DELETE("/users/:id", userHandler.Delete)

				admin.GET("/settings", settingHandler.GetAll)
				admin.POST("/settings", settingHandler.Set)
			}
		}
	}

	return r
}
