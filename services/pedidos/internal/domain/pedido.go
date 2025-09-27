package domain

import (
	"time"
)

// Status representa o estado atual de um pedido.
type Status string

// Os possíveis estados de um pedido
const (
	StatusAguardandoPagamento Status = "aguardando_pagamento"
	StatusPago                Status = "pago"
	StatusEnviado             Status = "enviado"
	StatusCancelado           Status = "cancelado"
)

// Item representa um item dentro do pedido.
type Item struct {
	ID         string
	ProdutoID  string
	Nome       string
	Preco      float64
	Quantidade int
}

// Pedido é a entidade raiz do nosso agregado.
type Pedido struct {
	ID           string
	ClienteID    string
	Itens        []Item
	Status       Status
	Total        float64
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

// NewPedido é o construtor do nosso agregado.
// Ele garante que o pedido seja criado de forma consistente.
func NewPedido(clienteID string, itens []Item) (*Pedido, error) {
	if len(itens) == 0 {
		return nil, ErrItemInvalido
	}

	total := 0.0
	for _, item := range itens {
		total += item.Preco * float64(item.Quantidade)
	}

	return &Pedido{
		ID:           "", // O ID será gerado na camada de infraestrutura
		ClienteID:    clienteID,
		Itens:        itens,
		Status:       StatusAguardandoPagamento,
		Total:        total,
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
	}, nil
}
