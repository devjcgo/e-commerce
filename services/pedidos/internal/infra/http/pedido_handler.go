package http

import (
	"database/sql"
	"ecommerce/pedidos/internal/application" // Verifique o import
	_ "ecommerce/pedidos/internal/domain"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// PedidoHandler lida com as requisições HTTP para pedidos.
type PedidoHandler struct {
	service *application.PedidoService
}

// NewPedidoHandler cria uma nova instância do handler de pedidos.
func NewPedidoHandler(service *application.PedidoService) *PedidoHandler {
	return &PedidoHandler{
		service: service,
	}
}

// requestBody define a estrutura esperada no corpo da requisição para criar um pedido.
type createRequestBody struct {
	ClienteID string                   `json:"cliente_id"`
	Itens     []application.ItensInput `json:"itens"`
}

// @Summary Cria um novo pedido
// @Description Cria um novo pedido com base nos dados do cliente e itens fornecidos.
// @Tags pedidos
// @Accept json
// @Produce json
// @Param pedido body createRequestBody true "Dados para criação do pedido"
// @Success 201
// @Failure 400 {string} string "Corpo da requisição inválido"
// @Failure 500 {string} string "Erro interno ao criar pedido"
// @Router /pedidos [post]
func (h *PedidoHandler) CriarPedidoHandler(w http.ResponseWriter, r *http.Request) {
	var body createRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	_, err := h.service.CriarPedido(r.Context(), body.ClienteID, body.Itens)
	if err != nil {
		http.Error(w, "Erro ao criar pedido: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Status 201 Created
	json.NewEncoder(w)
}

// @Summary Busca um pedido por ID
// @Description Retorna os detalhes de um pedido específico com base no seu UUID.
// @Tags pedidos
// @Produce json
// @Param id path string true "ID do Pedido (UUID)"
// @Success 200 {object} domain.Pedido
// @Failure 400 {string} string "O ID do pedido é obrigatório"
// @Failure 404 {string} string "Pedido não encontrado"
// @Failure 500 {string} string "Erro interno ao buscar pedido"
// @Router /pedidos/{id} [get]
func (h *PedidoHandler) BuscarPedidoPorIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extrai o 'id' do parâmetro da URL usando chi.
	pedidoID := chi.URLParam(r, "id")
	if pedidoID == "" {
		http.Error(w, "O ID do pedido é obrigatório", http.StatusBadRequest)
		return
	}

	// Chama o serviço da camada de aplicação.
	pedido, err := h.service.BuscarPedidoPorID(r.Context(), pedidoID)
	if err != nil {
		// Se o erro for 'sql.ErrNoRows', significa que não encontramos o pedido.
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Pedido não encontrado", http.StatusNotFound)
			return
		}
		// Para outros tipos de erro, retornamos um erro genérico.
		http.Error(w, "Erro ao buscar pedido: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	json.NewEncoder(w).Encode(pedido)
}

// @Summary Lista todos pedidos
// @Description Retorna pedidos e seus itens.
// @Tags pedidos
// @Produce json
// @Success 200 {object} []domain.Pedido
// @Failure 404 {string} string "Sem pedidos na base"
// @Failure 500 {string} string "Erro interno ao listar pedidos"
// @Router /pedidos [get]
func (h *PedidoHandler) ListarTodosPedidos(w http.ResponseWriter, r *http.Request) {

	pedidos, err := h.service.ListarPedidos(r.Context())
	if err != nil {

		// Se o erro for 'sql.ErrNoRows', significa que não encontramos o pedido.
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Sem pedidos na base", http.StatusNotFound)
			return
		}

		http.Error(w, "Erro ao listar pedidos "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	json.NewEncoder(w).Encode(pedidos)
}
