package domain

import "context"

// ClienteRepository define a interface para interações com o banco de dados de clientes.
type ClienteRepository interface {
	Save(ctx context.Context, cliente *Cliente) error
	FindAll(ctx context.Context) ([]*Cliente, error)
}
