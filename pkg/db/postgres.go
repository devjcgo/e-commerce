package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewConnection cria e retorna uma nova conexão com o banco de dados PostgreSQL.
func NewConnection() (*sql.DB, error) {
	// Pega a string de conexão da variável de ambiente injetada pelo Cloud Run.
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, fmt.Errorf("a variável de ambiente DATABASE_URL não foi definida")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão com o DB: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao pingar o DB: %w", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")
	return db, nil
}
