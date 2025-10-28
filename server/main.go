package main

import (
	"context"
	"log"
	"net"

	// Importa o pacote gRPC
	"google.golang.org/grpc"

	// Importa o pacote gerado
	pb "github.com/diegocamargo55555/RPC/proto"
)

// Define a porta em que o servidor escutará
const (
	port = ":50051"
)

// 'server' é usado para implementar a interface GreeterServer.
// Incorporamos 'UnimplementedGreeterServer' para compatibilidade futura.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implementa a função RPC definida no .proto
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// Logamos qual cliente fez a requisição.
	// Isso mostrará as chamadas de múltiplos clientes.
	log.Printf("Recebida requisição de: %v", in.GetName())

	// Retornamos a resposta
	return &pb.HelloReply{Message: "Olá, " + in.GetName() + "! Seja bem-vindo(a)."}, nil
}

func main() {
	// [REQUISITO B] Escutar em todas as interfaces de rede
	// Usar ":50051" (ou "0.0.0.0:50051") é crucial para aceitar conexões
	// de outras máquinas, e não apenas de 'localhost'.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Falha ao escutar na porta %s: %v", port, err)
	}

	log.Printf("Servidor gRPC escutando em %v", lis.Addr())

	// Cria uma nova instância do servidor gRPC
	s := grpc.NewServer()

	// Registra nosso serviço 'Greeter' no servidor gRPC
	pb.RegisterGreeterServer(s, &server{})

	// Inicia o servidor. Ele rodará indefinidamente, aceitando
	// múltiplas conexões de clientes [REQUISITO A]
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Falha ao iniciar servidor: %v", err)
	}
}
