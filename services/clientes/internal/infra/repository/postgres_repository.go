package repository

import (
	"context"
	"database/sql"
	"ecommerce/clientes/internal/domain"
	"time"

	"github.com/google/uuid"
)

type postgresClienteRepository struct {
	db *sql.DB
}

// NewPostgresClienteRepository é o construtor do nosso repositório.
func NewPostgresClienteRepository(db *sql.DB) domain.ClienteRepository {
	return &postgresClienteRepository{db: db}
}

// Save cria um novo cliente e seus endereços dentro de uma transação.
func (r *postgresClienteRepository) Save(ctx context.Context, cliente *domain.Cliente) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Gera um novo ID e define as datas de criação e alteração.
	cliente.ID = uuid.NewString()
	now := time.Now()
	cliente.CriadoEm = now
	cliente.AlteradoEm = now // <-- ALTERADO: Na criação, AlteradoEm é igual a CriadoEm.

	// Insere o cliente principal, incluindo a nova coluna.
	clienteQuery := `INSERT INTO clientes (id, nome, email, criado_em, alterado_em) VALUES ($1, $2, $3, $4, $5)`              // <-- ALTERADO
	_, err = tx.ExecContext(ctx, clienteQuery, cliente.ID, cliente.Nome, cliente.Email, cliente.CriadoEm, cliente.AlteradoEm) // <-- ALTERADO
	if err != nil {
		return err
	}

	// Insere os endereços associados (nenhuma mudança aqui)
	enderecoQuery := `INSERT INTO cliente_enderecos (cliente_id, rua, cidade, estado, cep) VALUES ($1, $2, $3, $4, $5)`
	for _, endereco := range cliente.Enderecos {
		_, err = tx.ExecContext(ctx, enderecoQuery, cliente.ID, endereco.Rua, endereco.Cidade, endereco.Estado, endereco.CEP)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FindAll busca todos os clientes e seus respectivos endereços.
func (r *postgresClienteRepository) FindAll(ctx context.Context) ([]*domain.Cliente, error) {
	// Adicionamos a coluna c.alterado_em à query.
	const query = `
		SELECT c.id, c.nome, c.email, c.criado_em, c.alterado_em, -- <-- ALTERADO
		       e.id, e.rua, e.cidade, e.estado, e.cep
		FROM clientes c
		LEFT JOIN cliente_enderecos e ON c.id = e.cliente_id
		ORDER BY c.criado_em DESC, c.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clientesMap := make(map[string]*domain.Cliente)
	var clientesOrdenados []*domain.Cliente

	for rows.Next() {
		var c domain.Cliente
		var e domain.Endereco
		var endID sql.NullInt64
		var endRua, endCidade, endEstado, endCEP sql.NullString

		// Adicionamos &c.AlteradoEm ao Scan.
		if err := rows.Scan(
			&c.ID, &c.Nome, &c.Email, &c.CriadoEm, &c.AlteradoEm, // <-- ALTERADO
			&endID, &endRua, &endCidade, &endEstado, &endCEP,
		); err != nil {
			return nil, err
		}

		if _, existe := clientesMap[c.ID]; !existe {
			c.Enderecos = []*domain.Endereco{}
			clientesMap[c.ID] = &c
			clientesOrdenados = append(clientesOrdenados, &c)
		}

		if endID.Valid {
			e.ID = endID.Int64
			e.Rua = endRua.String
			e.Cidade = endCidade.String
			e.Estado = endEstado.String
			e.CEP = endCEP.String
			clientesMap[c.ID].Enderecos = append(clientesMap[c.ID].Enderecos, &e)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clientesOrdenados, nil
}
