package repository

import (
	"context"
	"database/sql"

	// Import CORRETO do domain, usando o nome do módulo definido no go.mod
	"ecommerce/pedidos/internal/domain"
	"time"

	// Import do UUID
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresPedidoRepository struct {
	db *sql.DB
}

// O construtor recebe a conexão pronta
func NewPostgresPedidoRepository(db *sql.DB) domain.PedidoRepository {
	return &postgresPedidoRepository{db: db}
}

// Save persiste um pedido no banco de dados usando uma transação.
func (r *postgresPedidoRepository) Save(ctx context.Context, pedido *domain.Pedido) error {
	pedido.ID = uuid.NewString()
	pedido.AtualizadoEm = time.Now()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	pedidoQuery := `INSERT INTO pedidos (id, cliente_id, status, total, criado_em, atualizado_em)
					 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, pedidoQuery, pedido.ID, pedido.ClienteID, pedido.Status, pedido.Total, pedido.CriadoEm, pedido.AtualizadoEm)
	if err != nil {
		return err
	}

	itemQuery := `INSERT INTO pedido_itens (pedido_id, produto_id, nome_produto, preco, quantidade)
				  VALUES ($1, $2, $3, $4, $5)`
	for _, item := range pedido.Itens {
		_, err = tx.ExecContext(ctx, itemQuery, pedido.ID, item.ProdutoID, item.Nome, item.Preco, item.Quantidade)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FindByID busca um pedido e seus itens pelo ID.
func (r *postgresPedidoRepository) FindByID(ctx context.Context, id string) (*domain.Pedido, error) {
	query := `SELECT
				p.id, p.cliente_id, p.status, p.total, p.criado_em, p.atualizado_em,
				i.id, i.produto_id, i.nome_produto, i.preco, i.quantidade
			  FROM pedidos p
			  LEFT JOIN pedido_itens i ON p.id = i.pedido_id
			  WHERE p.id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pedido *domain.Pedido
	for rows.Next() {
		var itemID sql.NullInt64
		var item domain.Item

		if pedido == nil {
			pedido = &domain.Pedido{}
			err := rows.Scan(
				&pedido.ID, &pedido.ClienteID, &pedido.Status, &pedido.Total, &pedido.CriadoEm, &pedido.AtualizadoEm,
				&itemID, &item.ProdutoID, &item.Nome, &item.Preco, &item.Quantidade,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err := rows.Scan(
				&sql.RawBytes{}, &sql.RawBytes{}, &sql.RawBytes{}, &sql.RawBytes{}, &sql.RawBytes{}, &sql.RawBytes{},
				&itemID, &item.ProdutoID, &item.Nome, &item.Preco, &item.Quantidade,
			)
			if err != nil {
				return nil, err
			}
		}

		if itemID.Valid {
			pedido.Itens = append(pedido.Itens, item)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if pedido == nil {
		return nil, sql.ErrNoRows
	}

	return pedido, nil
}
