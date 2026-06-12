package main

import (
	"backend/clients"
	"backend/controllers"
	"backend/dao"
	"backend/services"
	"backend/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := clients.NewDBConnection()
	if err != nil {
		log.Fatal("error conectando a la base de datos:", err)
	}
	clients.AutoMigrate(db)

	// DAOs
	userDAO := dao.NewUserDAO(db)
	eventDAO := dao.NewEventDAO(db)
	ticketDAO := dao.NewTicketDAO(db)
	waitlistDAO := dao.NewWaitlistDAO(db)

	// Services
	userService := services.NewUserService(userDAO)
	eventService := services.NewEventService(eventDAO, ticketDAO)
	ticketService := services.NewTicketService(ticketDAO, eventDAO, userDAO, waitlistDAO)
	waitlistService := services.NewWaitlistService(waitlistDAO, eventDAO)

	// Controllers
	healthCtrl := controllers.NewHealthController(db)
	authCtrl := controllers.NewAuthController(userService)
	eventCtrl := controllers.NewEventController(eventService)
	ticketCtrl := controllers.NewTicketController(ticketService)
	waitlistCtrl := controllers.NewWaitlistController(waitlistService)

	r := gin.Default()
	r.Use(utils.CORSMiddleware())
	r.GET("/health", healthCtrl.Check)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}

	// Eventos — públicos
	r.GET("/events", eventCtrl.List)
	r.GET("/events/:id", eventCtrl.Get)

	// Rutas protegidas (requieren JWT)
	protected := r.Group("/")
	protected.Use(utils.AuthMiddleware())
	{
		protected.POST("/tickets", ticketCtrl.Buy)
		protected.GET("/tickets/mine", ticketCtrl.GetMine)
		protected.DELETE("/tickets/:id", ticketCtrl.Cancel)
		protected.PUT("/tickets/:id/transfer", ticketCtrl.Transfer)

		protected.POST("/events/:id/waitlist", waitlistCtrl.Join)
		protected.GET("/waitlist/mine", waitlistCtrl.GetMine)
	}

	// Rutas de administrador (requieren JWT + rol admin)
	admin := r.Group("/")
	admin.Use(utils.AuthMiddleware(), utils.AdminMiddleware())
	{
		admin.GET("/admin/events", eventCtrl.ListAll)
		admin.POST("/events", eventCtrl.Create)
		admin.PUT("/events/:id", eventCtrl.Update)
		admin.PATCH("/events/:id/cancel", eventCtrl.Cancel)
		admin.GET("/events/:id/report", eventCtrl.Report)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
