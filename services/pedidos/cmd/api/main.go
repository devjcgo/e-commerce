package main

import (
	"ecommerce/pedidos/internal/application"
	httphandler "ecommerce/pedidos/internal/infrastructure/http"
	"ecommerce/pedidos/internal/infrastructure/repository"
	"ecommerce/pkg/db"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "ecommerce/pedidos/docs" // Importa os docs gerados pelo swag (necessário)

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger" // Importa o handler do swagger
)

// @title API de Pedidos do E-commerce
// @version 1.0
// @description Este é o microsserviço responsável pelo gerenciamento de pedidos.
// @host localhost:8080
// @BasePath /
func main() {
	// 1. Inicializa a Conexão com o Banco de Dados
	dbConn, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer dbConn.Close()

	// 2. Inicializa o Repositório, Serviço e Handler
	repo := repository.NewPostgresPedidoRepository(dbConn)
	pedidoService := application.NewPedidoService(repo)
	pedidoHandler := httphandler.NewPedidoHandler(pedidoService)

	// 3. Configuração do Roteador e Rotas
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rota para a API
	r.Post("/pedidos", pedidoHandler.CriarPedidoHandler)
	r.Get("/pedidos/{id}", pedidoHandler.BuscarPedidoPorIDHandler)

	// Rota para a documentação do Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// 4. Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		// Se não estiver no Cloud Run, usa 8080 como padrão para rodar localmente.
		port = "8080"
	}

	// Monta o endereço de escuta, ex: ":8080"
	addr := fmt.Sprintf(":%s", port)

	fmt.Printf("Servidor de Pedidos rodando na porta %s...\n", port)
	fmt.Printf("Acesse a documentação da API em http://localhost:%s/swagger/index.html\n", port)

	// Usa a variável 'addr' para iniciar o servidor.
	log.Fatal(http.ListenAndServe(addr, r))
}
