// server/main.go
package main

import (
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "stock-app/stockpb"
)

type server struct {
	pb.UnimplementedStockQuoteServiceServer
}

func (s *server) GetStockQuotes(req *pb.QuoteRequest, stream pb.StockQuoteService_GetStockQuotesServer) error {
	ticker := req.GetTicker()
	log.Printf("Cliente solicitou cotações para: %s", ticker)

	currentPrice := rand.Float64()*400.0 + 100.0

	for {
		if err := stream.Context().Err(); err != nil {
			log.Printf("Cliente para %s desconectado: %v", ticker, err)
			return status.Errorf(codes.Canceled, "Cliente desconectado")
		}

		change := rand.Float64()*2.0 - 1.0
		currentPrice += change

		if currentPrice < 0.1 {
			currentPrice = 0.1
		}

		response := &pb.QuoteResponse{
			Ticker:    ticker,
			Price:     math.Round(currentPrice*100) / 100,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(response); err != nil {
			log.Printf("Erro ao enviar stream para %s: %v", ticker, err)
			return err
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	log.Println("Iniciando servidor gRPC na porta 50051...")

	lis, err := net.Listen("tcp", "[::]:50051")
	if err != nil {
		log.Fatalf("Falha ao ouvir: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterStockQuoteServiceServer(s, &server{})

	log.Println("Servidor pronto e aguardando conexões.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}
