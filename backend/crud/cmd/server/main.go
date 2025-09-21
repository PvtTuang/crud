package main

import (
	"auth/middleware"
	"auth/pkg/authentication"
	"crud/pkg/cart"
	"crud/pkg/history"
	"crud/pkg/product"

	"database/config"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()

	// migrate models
	config.Migrate(db,
		&product.Product{},
		&cart.Cart{},
		&cart.CartProduct{},
		&history.PurchaseHistory{},
		&history.PurchaseItem{},
	)

	// Product setup
	productRepo := product.NewProductRepository(db)
	productService := product.NewProductService(productRepo)
	productHandler := product.NewProductHandler(productService)

	// Cart setup
	cartRepo := cart.NewCartRepository(db)
	cartService := cart.NewCartService(cartRepo)
	cartHandler := cart.NewCartHandler(cartService)

	// History setup
	historyRepo := history.NewHistoryRepository(db)
	historyService := history.NewHistoryService(historyRepo)
	historyHandler := history.NewHistoryHandler(historyService)

	// Gin setup
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

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "CRUD Service OK"})
	})

	// Auth service (JWT)
	rdb := config.ConnectRedis()
	userRepo := authentication.NewUserRepository(db)
	authService := authentication.NewAuthService(userRepo, rdb, []byte(os.Getenv("JWT_SECRET")))

	// Product routes
	products := r.Group("/products")
	products.Use(middleware.JWTMiddleware(authService))
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.GetAllProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Cart routes
	carts := r.Group("/carts")
	carts.Use(middleware.JWTMiddleware(authService))
	{
		carts.GET("/:user_id", cartHandler.GetCart)
		carts.POST("/:user_id/products", cartHandler.AddProduct)
		carts.DELETE("/:user_id/products/:product_id", cartHandler.RemoveProduct)
		carts.DELETE("/:user_id/clear", cartHandler.ClearCart)

		carts.POST("/:user_id/checkout", cartHandler.Checkout)
	}

	// History routes
	histories := r.Group("/history")
	histories.Use(middleware.JWTMiddleware(authService))
	{
		histories.GET("/:user_id", historyHandler.GetByUser)
	}

	// Run server
	port := os.Getenv("CRUD_HTTP_PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("CRUD Service รันที่ port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("ไม่สามารถรันเซิร์ฟเวอร์:", err)
	}
}
