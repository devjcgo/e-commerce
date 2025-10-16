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
		var item *domain.Item

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

func (r *postgresPedidoRepository) ListAll(ctx context.Context) ([]*domain.Pedido, error) {
	// 1. A QUERY: Busca todos os dados de uma vez, ordenando pelos pedidos
	// para garantir que as linhas do mesmo pedido venham em sequência.
	const query = `
		SELECT
			p.id, p.cliente_id, p.status, p.total, p.criado_em, p.atualizado_em,
			i.id, i.produto_id, i.nome_produto, i.preco, i.quantidade
		FROM pedidos p
		LEFT JOIN pedido_itens i ON p.id = i.pedido_id
		ORDER BY p.criado_em DESC, p.id` // Ordenação estável

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. ESTRUTURAS DE APOIO:
	// - O 'map' para evitar duplicar pedidos e agrupar itens rapidamente.
	// - O 'slice' para manter a ordem original que veio do banco de dados.
	pedidosMap := make(map[string]*domain.Pedido)
	var pedidosOrdenados []*domain.Pedido

	// 3. O LOOP DE PROCESSAMENTO
	for rows.Next() {
		var p domain.Pedido
		var item domain.Item
		// Usamos tipos que aceitam NULL para as colunas de 'pedido_itens',
		// pois um pedido pode não ter itens.
		var itemID sql.NullInt64
		var itemProdutoID sql.NullString
		var itemNome sql.NullString
		var itemPreco sql.NullFloat64
		var itemQuantidade sql.NullInt32

		if err := rows.Scan(
			&p.ID, &p.ClienteID, &p.Status, &p.Total, &p.CriadoEm, &p.AtualizadoEm,
			&itemID, &itemProdutoID, &itemNome, &itemPreco, &itemQuantidade,
		); err != nil {
			return nil, err
		}

		// 4. LÓGICA DE AGRUPAMENTO
		// Se o pedido ainda não está no nosso map...
		if _, existe := pedidosMap[p.ID]; !existe {
			// ...é um novo pedido. Inicializamos sua lista de itens...
			p.Itens = []*domain.Item{}
			// ...adicionamos ao map para encontrá-lo nas próximas linhas...
			pedidosMap[p.ID] = &p
			// ...e adicionamos ao slice para preservar a ordem.
			pedidosOrdenados = append(pedidosOrdenados, &p)
		}

		// 5. ADIÇÃO DO ITEM
		// Se esta linha contém um item válido (itemID não é NULL)...
		if itemID.Valid {
			// ...criamos o struct do item...
			item.ID = string(itemID.Int64)
			item.ProdutoID = itemProdutoID.String
			item.Nome = itemNome.String
			item.Preco = itemPreco.Float64
			item.Quantidade = int(itemQuantidade.Int32)

			// ...e o adicionamos à lista de itens do pedido correto (que buscamos no map).
			pedidosMap[p.ID].Itens = append(pedidosMap[p.ID].Itens, &item)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 6. RETORNO
	return pedidosOrdenados, nil
}
