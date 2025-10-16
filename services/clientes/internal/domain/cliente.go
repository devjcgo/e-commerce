package domain

import "time"

// Cliente é a nossa raiz de agregado.
type Cliente struct {
	ID         string
	Nome       string
	Email      string
	Enderecos  []*Endereco
	CriadoEm   time.Time
	AlteradoEm time.Time
}

// Endereco pertence ao agregado de Cliente.
type Endereco struct {
	ID     int64
	Rua    string
	Cidade string
	Estado string
	CEP    string
}
