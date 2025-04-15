package services

import (
	"context"

	"github.com/YuMingChang/liquidity-provider.git/internal/market"
	"github.com/YuMingChang/liquidity-provider.git/internal/models"
	"github.com/YuMingChang/liquidity-provider.git/internal/repositories"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderBookService struct {
	repo       *repositories.OrderRepository
	marketConn *grpc.ClientConn
}

func NewOrderBookService(repo *repositories.OrderRepository) *OrderBookService {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &OrderBookService{repo: repo, marketConn: conn}
}

func (s *OrderBookService) PlaceOrder(symbol string, side string, price float64, quantity float64) error {
	order := &models.Order{
		Symbol:   symbol,
		Price:    price,
		Quantity: quantity,
		Side:     side,
		Status:   "open",
	}
	return s.repo.CreateOrder(order)
}

func (s *OrderBookService) MatchOrders(symbol string) error {
	orders, err := s.repo.GetOpenOrders(symbol)
	if err != nil {
		return err
	}

	for _, buy := range orders {
		if buy.Side != "buy" || buy.Status != "open" {
			continue
		}
		for _, sell := range orders {
			if sell.Side != "sell" || sell.Status != "open" {
				continue
			}
			if buy.Price >= sell.Price {
				buy.Status = "closed"
				sell.Status = "closed"
				s.repo.UpdateOrder(&buy)
				s.repo.UpdateOrder(&sell)
				break
			}
		}
	}
	return nil
}

func (s *OrderBookService) GetMarketData(symbol string) (*market.MarketDataResponse, error) {
	client := market.NewMarketDataServiceClient(s.marketConn)
	return client.GetMarketData(context.Background(), &market.MarketDataRequest{Symbol: symbol})
}
