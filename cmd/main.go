package cmd

import (
	"fmt"
	"go-rental/internal/customer"
	"go-rental/internal/rent"
	"go-rental/internal/user"
	"go-rental/internal/vehicle"
	"go-rental/pkg/config"
	"go-rental/pkg/middlewares"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %d - %s %s %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.StatusCode,
			param.Method,
			param.Path,
			param.Latency,
		)
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CorsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	r.Use(middlewares.GinErrorHandler())
	
	// === Database ===
	if err := config.Connect(cfg); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	db := config.GetDB()
	tables := []interface{}{
		&user.User{},
		&vehicle.Vehicle{},
		&customer.Customer{},
		&rent.Rent{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("✅ Migrasi database berhasil.")

	// === Home Route ===
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "Welcome to the REST API",
			"version":   "1.0.0",
			"timestamp": time.Now(),
		})
	})

	// === Services, Repositories & Controllers ===
	emailService := email.NewService(cfg)

	userRepo := user.NewRepository(db)
	eventRepo := event.Newrepository(db)
	participantRepo := participant.Newrepository(db)
	scheduleRepo := schedule.NewRepository(db)
	notificationRepo := notification.Newrepository(db)

	eventRepoAdapter := event.NewEventRepositoryAdapter(eventRepo)

	notificationService := notification.NewService(notificationRepo, eventRepo, emailService, cfg)
	notificationController := notification.NewController(notificationService, cfg)

	userService := user.NewService(userRepo, emailService, cfg)
	userController := user.NewController(userService, cfg)

	eventService := event.NewService(eventRepo, participantRepo, userRepo, notificationService, cfg)
	eventController := event.NewController(eventService, cfg)

	participantService := participant.NewService(participantRepo, eventRepoAdapter, userRepo, emailService, cfg)
	participantController := participant.NewController(participantService, *cfg)

	scheduleService := schedule.NewService(scheduleRepo, eventRepo, cfg)
	scheduleController := schedule.NewController(scheduleService, cfg)

	// === Scheduler ===
	scheduler := schedule.NewScheduler(scheduleRepo, notificationService, participantRepo, userRepo)
	scheduler.Start()
	defer scheduler.Stop()

	// === Routes Setup (converted to Gin) ===
	user.SetupUserRoutesGin(r, userController, cfg)
	event.SetupOrganizerEventRoutesGin(r, eventController, cfg)
	participant.SetupParticipantRoutesGin(r, participantController, cfg)
	schedule.SetupScheduleRoutesGin(r, scheduleController, cfg)
	notification.SetupNotificationRoutesGin(r, notificationController, cfg)

	// 404 Not Found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Route not found"})
	})

	// === Start Server ===
	log.Printf("Server running on port %s", cfg.Port)
	log.Printf("Local: http://localhost:%s", cfg.Port)
	log.Printf("Environment: %s", cfg.NodeEnv)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}