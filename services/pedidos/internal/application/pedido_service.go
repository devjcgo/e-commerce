package application

import (
	"context"
	"ecommerce/pedidos/internal/domain" // Verifique se o import está correto
)

// PedidoService é a implementação dos nossos casos de uso de pedido.
type PedidoService struct {
	repo domain.PedidoRepository
}

// NewPedidoService é o construtor do nosso serviço de aplicação.
func NewPedidoService(repo domain.PedidoRepository) *PedidoService {
	return &PedidoService{
		repo: repo,
	}
}

// ItensInput é um DTO para os itens na criação do pedido.
type ItensInput struct {
	ProdutoID  string  `json:"produto_id"`
	Nome       string  `json:"nome"`
	Preco      float64 `json:"preco"`
	Quantidade int     `json:"quantidade"`
}

// CriarPedido é o caso de uso para criar um novo pedido.
func (s *PedidoService) CriarPedido(ctx context.Context, clienteID string, itensInput []ItensInput) (*domain.Pedido, error) {
	var itensDominio []*domain.Item
	for _, itemInput := range itensInput {
		itensDominio = append(itensDominio, &domain.Item{
			ProdutoID:  itemInput.ProdutoID,
			Nome:       itemInput.Nome,
			Preco:      itemInput.Preco,
			Quantidade: itemInput.Quantidade,
		})
	}

	novoPedido, err := domain.NewPedido(clienteID, itensDominio)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, novoPedido)
	if err != nil {
		return nil, err
	}

	return novoPedido, nil
}

func (s *PedidoService) BuscarPedidoPorID(ctx context.Context, id string) (*domain.Pedido, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *PedidoService) ListarPedidos(ctx context.Context) ([]*domain.Pedido, error) {
	return s.repo.ListAll(ctx)
}
