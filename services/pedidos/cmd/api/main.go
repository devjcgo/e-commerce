package main

import (
	"ecommerce/pedidos/internal/application"
	httphandler "ecommerce/pedidos/internal/infra/http"
	"ecommerce/pedidos/internal/infra/repository"
	"ecommerce/pkg/db"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "ecommerce/pedidos/docs" // Importa os docs gerados pelo swag (necessário)

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger" // Importa o handler do swagger
)

// @title API de Pedidos do E-commerce
// @version 1.0
// @description Este é o microsserviço responsável pelo gerenciamento de pedidos.
// @BasePath /pedidos
func main() {
	// --- CARREGA O ARQUIVO .ENV (APENAS PARA AMBIENTE LOCAL) ---
	err := godotenv.Load()
	if err != nil {
		// Não usamos log.Fatalf porque em produção o arquivo .env não existirá
		// e não queremos que a aplicação pare. Apenas avisamos no console.
		log.Println("Aviso: Erro ao carregar arquivo .env")
	}

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

	// Rotas da API
	r.Post("/pedidos", pedidoHandler.CriarPedidoHandler)
	r.Get("/pedidos/{id}", pedidoHandler.BuscarPedidoPorIDHandler)
	r.Get("/pedidos", pedidoHandler.ListarTodosPedidos)

	// Rota para a documentação do Swagger (AGORA CORRIGIDA)
	r.Get("/swagger/*", httpSwagger.Handler())

	// 4. Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)

	fmt.Printf("Servidor de Pedidos rodando na porta %s...\n", port)
	fmt.Printf("Acesse a documentação da API em http://localhost:%s/swagger/index.html\n", port)

	log.Fatal(http.ListenAndServe(addr, r))
}
