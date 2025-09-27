package domain

import "context"

// PedidoRepository define os métodos para persistir e recuperar pedidos.
type PedidoRepository interface {
	Save(ctx context.Context, pedido *Pedido) error
	FindByID(ctx context.Context, id string) (*Pedido, error)
	// Outros métodos de consulta, como FindAll, etc.
}
