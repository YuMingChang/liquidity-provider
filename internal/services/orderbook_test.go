package services

import (
	"testing"

	"github.com/YuMingChang/liquidity-provider.git/internal/models"
	"github.com/YuMingChang/liquidity-provider.git/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPlaceOrder(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Order{})
	repo := repositories.NewOrderRepository(db)
	service := NewOrderBookService(repo)

	err := service.PlaceOrder("BTC/USD", "buy", 50000, 0.1)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}

	orders, _ := repo.GetOpenOrders("BTC/USD")
	if len(orders) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(orders))
	}
}

func TestMatchOrders(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Order{})
	repo := repositories.NewOrderRepository(db)
	service := NewOrderBookService(repo)

	service.PlaceOrder("BTC/USD", "buy", 50000, 0.1)
	service.PlaceOrder("BTC/USD", "sell", 49000, 0.1)
	service.MatchOrders("BTC/USD")

	orders, _ := repo.GetOpenOrders("BTC/USD")
	if len(orders) != 0 {
		t.Fatalf("Expected 0 open orders after matching, got %d", len(orders))
	}
}

func TestApplyGridStrategy(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Order{})
	repo := repositories.NewOrderRepository(db)
	service := NewOrderBookService(repo)

	err := service.ApplyGridStrategy("BTC/USD", 100, 2)
	if err != nil {
		t.Fatalf("Failed to apply grid strategy: %v", err)
	}

	orders, _ := repo.GetOpenOrders("BTC/USD")
	if len(orders) != 4 { // 2 buy and 2 sell orders
		t.Fatalf("Expected 4 open orders, got %d", len(orders))
	}
}
