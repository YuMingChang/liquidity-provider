package market

import (
	"context"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
)

type MarketDataServer struct {
	UnimplementedMarketDataServiceServer
}

func (s *MarketDataServer) GetMarketData(ctx context.Context, req *MarketDataRequest) (*MarketDataResponse, error) {
	// Simulate market data
	rand.New(rand.NewSource(time.Now().UnixNano()))
	price := 50000.0 + rand.Float64()*1000.0 // Simulated BTC/USD price
	volume := 10.0 + rand.Float64()*5.0      // Simulated volume
	return &MarketDataResponse{
		Symbol: req.Symbol,
		Price:  price,
		Volume: volume,
	}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	RegisterMarketDataServiceServer(s, &MarketDataServer{})
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
