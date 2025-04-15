package repositories

import (
	"github.com/YuMingChang/liquidity-provider.git/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetOpenOrders(symbol string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("symbol = ? AND status = ?", symbol, "open").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	return r.db.Save(order).Error
}
