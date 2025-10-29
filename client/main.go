// client/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Importe o pacote gerado
	pb "stock-app/stockpb"
)

func runClient(ticker, serverAddress string) {
	target := fmt.Sprintf("%s:50051", serverAddress)
	log.Printf("Tentando conectar ao servidor em %s para o ticker %s...", target, ticker)

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewStockQuoteServiceClient(conn)

	request := &pb.QuoteRequest{
		Ticker: strings.ToUpper(ticker),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := client.GetStockQuotes(ctx, request)
	if err != nil {
		log.Fatalf("Erro ao chamar GetStockQuotes: %v", err)
	}

	log.Printf("--- Recebendo cotações para %s ---", ticker)

	for {
		response, err := stream.Recv()

		if err == io.EOF {
			log.Println("Stream finalizado pelo servidor.")
			break
		}
		if err != nil {
			log.Fatalf("Erro ao receber stream: %v", err)
		}

		log.Printf("[Ticker: %s] Preço: $%.2f (Timestamp: %d)",
			response.GetTicker(),
			response.GetPrice(),
			response.GetTimestamp())
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter o ticket da acao: ")
	tickerInput, _ := reader.ReadString('\n')
	ticker := strings.TrimSpace(tickerInput)
	serverAddress := "192.168.18.243"

	runClient(ticker, serverAddress)
}
