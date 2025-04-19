package handlers

import (
	"net/http"

	"github.com/YuMingChang/liquidity-provider.git/internal/services"
	"github.com/gin-gonic/gin"
)

type OrderBookHandler struct {
	service *services.OrderBookService
}

func NewOrderBookHandler(service *services.OrderBookService) *OrderBookHandler {
	return &OrderBookHandler{service: service}
}

func (h *OrderBookHandler) PlaceOrder(c *gin.Context) {
	var req struct {
		Symbol   string  `json:"symbol"`
		Side     string  `json:"side"`
		Price    float64 `json:"price"`
		Quantity float64 `json:"quantity"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check slippage and position limits
	if ok, err := h.service.CheckSlippage(req.Symbol, req.Price); err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slippage too high"})
		return
	}
	if ok, err := h.service.CheckPositionLimit(req.Symbol); err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Position limit exceeded"})
		return
	}

	if err := h.service.PlaceOrder(req.Symbol, req.Side, req.Price, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order placed"})
}

func (h *OrderBookHandler) ApplyGridStrategy(c *gin.Context) {
	var req struct {
		Symbol   string  `json:"symbol"`
		GridSize float64 `json:"grid_size"`
		Levels   int     `json:"levels"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ApplyGridStrategy(req.Symbol, req.GridSize, req.Levels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Grid strategy applied"})
}
