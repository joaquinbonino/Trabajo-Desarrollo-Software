package main

import (
	"backend/clients"
	"backend/controllers"
	"backend/dao"
	"backend/services"
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

	// Services
	userService := services.NewUserService(userDAO)
	_ = services.NewEventService(eventDAO)
	_ = services.NewTicketService(ticketDAO, eventDAO, userDAO)

	// Controllers
	healthCtrl := controllers.NewHealthController(db)
	authCtrl := controllers.NewAuthController(userService)

	r := gin.Default()
	r.GET("/health", healthCtrl.Check)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
