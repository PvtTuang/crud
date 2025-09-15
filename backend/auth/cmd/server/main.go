package main

import (
	"auth/pkg/authentication"
	"database/config"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET ต้องถูกกำหนดใน .env")
	}

	db := config.ConnectDB()
	config.Migrate(db, &authentication.User{})
	rdb := config.ConnectRedis()

	userRepo := authentication.NewUserRepository(db)
	authService := authentication.NewAuthService(userRepo, rdb, []byte(jwtSecret))
	authHandler := authentication.NewAuthHandler(authService)

	//Gin
	r := gin.Default()

	reactPort := os.Getenv("REACT_APP_API_URL")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{reactPort},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Auth Service OK"})
	})
	r.HEAD("/health", func(ctx *gin.Context) {
		ctx.Status(200)
	})
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	port := os.Getenv("AUTH_HTTP_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Auth service running at :%s", port)
	r.Run(":" + port)
}
