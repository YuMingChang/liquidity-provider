package main

import (
	"github.com/YuMingChang/liquidity-provider.git/internal/handlers"
	"github.com/YuMingChang/liquidity-provider.git/internal/market"
	"github.com/YuMingChang/liquidity-provider.git/internal/models"
	"github.com/YuMingChang/liquidity-provider.git/internal/repositories"
	"github.com/YuMingChang/liquidity-provider.git/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Start the gRPC server for market data
	go market.StartGRPCServer()

	// Initialize database
	db, err := gorm.Open(sqlite.Open("orders.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Order{})

	// Initialize layers
	orderRepo := repositories.NewOrderRepository(db)
	orderBookService := services.NewOrderBookService(orderRepo)
	orderBookHandler := handlers.NewOrderBookHandler(orderBookService)

	// Set up GIN router
	r := gin.Default()
	r.POST("/orders", orderBookHandler.PlaceOrder)
	r.POST("/grid", orderBookHandler.ApplyGridStrategy)

	// Start HTTP server
	r.Run(":8080")
}
