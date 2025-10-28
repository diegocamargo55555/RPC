package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/diegocamargo55555/RPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // Importa o pacote 'insecure'
)

// [REQUISITO B] Endereço do servidor.
// !! IMPORTANTE !! Mude "localhost" para o IP da máquina do servidor.
const (
	defaultAddress = "localhost:50051"
	defaultName    = "Mundo"
)

func main() {
	// Endereço do servidor. Use o IP da máquina servidora se estiver em outra máquina.
	// Ex: serverAddr := "192.168.1.10:50051"
	serverAddr := defaultAddress

	// [REQUISITO A] Pega um nome da linha de comando para diferenciar os clientes
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// [REQUISITO B] Configurando a conexão.
	// Usamos 'insecure.NewCredentials()' para este exemplo, pois não
	// configuramos TLS/SSL. Em produção, você DEVE usar credenciais seguras.
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close() // Fecha a conexão ao final

	log.Printf("Conectado ao servidor em %s", serverAddr)

	// Cria o "stub" do cliente
	c := pb.NewGreeterClient(conn)

	// Define um timeout para a chamada RPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Chama a função remota 'SayHello'
	log.Printf("Enviando saudação para o servidor com o nome: %s", name)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Não foi possível saudar: %v", err)
	}

	// Imprime a resposta do servidor
	log.Printf("Resposta do Servidor: %s", r.GetMessage())
}
