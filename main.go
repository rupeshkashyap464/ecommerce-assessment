package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := InitDB("db.json")
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	r := gin.Default()

	// --- CORS middleware so React (localhost:3000) can call Go (localhost:8080) ---
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// ---------------------------------------------------------------------------

	// users
	r.POST("/users", CreateUser)
	r.GET("/users", ListUsers)
	r.POST("/users/login", LoginUser)

	// items
	r.POST("/items", CreateItem)
	r.GET("/items", ListItems)

	// carts (token required)
	carts := r.Group("/carts")
	carts.Use(AuthRequired())
	carts.POST("/", CreateOrAddToCart)
	r.GET("/carts", ListCarts) // public list as per assignment

	// orders (token required)
	orders := r.Group("/orders")
	orders.Use(AuthRequired())
	orders.POST("/", CreateOrderFromCart)
	r.GET("/orders", ListOrders)

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("listening on :8080")
	r.Run(":8080")
}
