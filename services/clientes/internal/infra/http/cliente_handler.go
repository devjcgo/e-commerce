package http

import (
	"ecommerce/clientes/internal/application"
	"encoding/json"
	"net/http"
)

// ClienteHandler lida com as requisições HTTP para clientes.
type ClienteHandler struct {
	service *application.ClienteService
}

// NewClienteHandler é o construtor do nosso handler.
func NewClienteHandler(service *application.ClienteService) *ClienteHandler {
	return &ClienteHandler{
		service: service,
	}
}

// @Summary Cria um novo cliente
// @Description Cria um novo cliente com seus dados e endereços.
// @Tags clientes
// @Accept json
// @Produce json
// @Param cliente body application.ClienteInput true "Dados para criação do cliente"
// @Success 201 {object} domain.Cliente
// @Failure 400 {string} string "Corpo da requisição inválido"
// @Failure 500 {string} string "Erro interno ao criar cliente"
// @Router /clientes [post]
func (h *ClienteHandler) CriarClienteHandler(w http.ResponseWriter, r *http.Request) {
	// Decodifica o corpo da requisição JSON para o nosso DTO de entrada.
	var input application.ClienteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	// Chama o serviço da camada de aplicação com os dados recebidos.
	cliente, err := h.service.CriarCliente(r.Context(), input)
	if err != nil {
		http.Error(w, "Erro ao criar cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Se tudo deu certo, codifica o cliente criado como JSON e envia na resposta.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Status 201 Created
	json.NewEncoder(w).Encode(cliente)
}

// @Summary Lista todos os clientes
// @Description Retorna uma lista de todos os clientes cadastrados com seus endereços.
// @Tags clientes
// @Produce json
// @Success 200 {array} domain.Cliente
// @Failure 500 {string} string "Erro interno ao listar clientes"
// @Router /clientes [get]
func (h *ClienteHandler) ListarClientesHandler(w http.ResponseWriter, r *http.Request) {
	// Chama o serviço da camada de aplicação.
	clientes, err := h.service.ListarClientes(r.Context())
	if err != nil {
		http.Error(w, "Erro ao listar clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Codifica o slice de clientes como JSON e envia na resposta.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	json.NewEncoder(w).Encode(clientes)
}
