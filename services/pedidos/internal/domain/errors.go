package domain

import "errors"

// Erros que podem ser retornados pela camada de domínio.
var (
	ErrPedidoNaoEncontrado = errors.New("pedido não encontrado")
	ErrStatusInvalido      = errors.New("status do pedido inválido")
	ErrItemInvalido        = errors.New("item do pedido inválido")
)
