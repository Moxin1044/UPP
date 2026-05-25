package main

import (
	"log"

	"upp/internal/api"
	"upp/internal/config"
	"upp/internal/model"
	"upp/internal/notification"
	"upp/internal/repository"
	"upp/internal/scheduler"
	"upp/internal/service"
	"upp/internal/websocket"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init PostgreSQL
	db, err := model.InitDB(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Init MongoDB
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("Failed to connect MongoDB: %v", err)
	}
	_ = mongoClient // reserved for future use
	mongoDB := mongoClient.Database(cfg.MongoDB.Database)
	_ = mongoDB

	// Init Repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	monitorRepo := repository.NewMonitorRepository(db)
	resultRepo := repository.NewMonitorResultRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	// Init Notification Channel & Dispatcher
	notifyChan := make(chan scheduler.NotificationEvent, 100)
	dispatcher := notification.NewDispatcher(notificationRepo)
	dispatcher.Start(notifyChan)

	// Init Scheduler
	sched := scheduler.NewScheduler(monitorRepo, resultRepo, notifyChan)
	sched.Start()

	// Init Services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.Expire)
	projectService := service.NewProjectService(projectRepo)
	monitorService := service.NewMonitorService(monitorRepo, sched)
	resultService := service.NewMonitorResultService(resultRepo)
	userService := service.NewUserService(userRepo)
	settingService := service.NewSettingService(settingRepo)
	notificationService := service.NewNotificationService(notificationRepo, dispatcher)

	// Init default admin
	if err := authService.InitAdmin(); err != nil {
		log.Printf("Warning: Failed to init admin: %v", err)
	}

	// Init WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// Init Handlers
	authHandler := api.NewAuthHandler(authService)
	projectHandler := api.NewProjectHandler(projectService)
	monitorHandler := api.NewMonitorHandler(monitorService, resultService)
	notificationHandler := api.NewNotificationHandler(notificationService)
	userHandler := api.NewUserHandler(userService)
	settingHandler := api.NewSettingHandler(settingService)
	dashboardHandler := api.NewDashboardHandler(monitorService, resultRepo)

	// Setup Router
	router := api.SetupRouter(
		authHandler,
		projectHandler,
		monitorHandler,
		notificationHandler,
		userHandler,
		settingHandler,
		dashboardHandler,
		hub,
		cfg.JWT.Secret,
	)

	log.Printf("UPP Server starting on :%s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
