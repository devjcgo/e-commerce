package main

import (
	"ecommerce/clientes/internal/application"
	httphandler "ecommerce/clientes/internal/infra/http"
	"ecommerce/clientes/internal/infra/repository"
	"ecommerce/pkg/db"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "ecommerce/clientes/docs" // <-- IMPORT DOS DOCS GERADOS

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger" // <-- IMPORT DO HANDLER DO SWAGGER
)

// @title API de Clientes do E-commerce
// @version 1.0
// @description Microsserviço responsável pelo gerenciamento de clientes.
// @BasePath /clientes
func main() {
	//
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Erro ao carregar arquivo .env")
	}

	dbConn, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer dbConn.Close()

	repo := repository.NewPostgresClienteRepository(dbConn)
	clienteService := application.NewClienteService(repo)
	clienteHandler := httphandler.NewClienteHandler(clienteService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/clientes", clienteHandler.CriarClienteHandler)
	r.Get("/clientes", clienteHandler.ListarClientesHandler)

	// --- ROTA DO SWAGGER ADICIONADA ---
	r.Get("/swagger/*", httpSwagger.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Servidor de Clientes rodando na porta %s...\n", port)
	fmt.Printf("Acesse a documentação da API em http://localhost:%s/swagger/index.html\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
