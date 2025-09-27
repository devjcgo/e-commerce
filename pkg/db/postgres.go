package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewConnection cria e retorna uma nova conexão com o banco de dados PostgreSQL.
func NewConnection() (*sql.DB, error) {
	host := getEnv("DB_HOST", "34.70.240.93")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "@K953zhok")
	dbname := getEnv("DB_NAME", "pedidos")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
