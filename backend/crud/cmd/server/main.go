package main

import (
	"auth/middleware"
	"auth/pkg/authentication"
	"crud/pkg"
	"database/config"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()

	config.Migrate(db, &pkg.Product{})

	productRepo := pkg.NewProductRepository(db)
	productService := pkg.NewProductService(productRepo)
	productHandler := pkg.NewProductHandler(productService)

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
		c.JSON(200, gin.H{"status": "CRUD Service OK"})
	})

	rdb := config.ConnectRedis()

	userRepo := authentication.NewUserRepository(db)
	authService := authentication.NewAuthService(userRepo, rdb, []byte(os.Getenv("JWT_SECRET")))

	products := r.Group("/products")
	products.Use(middleware.JWTMiddleware(authService))
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.GetAllProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	port := os.Getenv("CRUD_HTTP_PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("CRUD Service รันที่ port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("ไม่สามารถรันเซิร์ฟเวอร์:", err)
	}
}
