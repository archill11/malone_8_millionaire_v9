package pg

import (
	"context"
	_ "embed"
	"fmt"
	"myapp/pkg/logger"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed schemes/users.sql
var users string

//go:embed schemes/schemes.sql
var schemes string

type (
	DBConfig struct {
		User     string
		Password string
		Database string
		Host     string
		Port     string
	}

	Database struct {
		db *pgxpool.Pool
		l  *logger.Logger
	}
)

func New(config DBConfig, l *logger.Logger) (*Database, error) {
	databaseURI := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, config.Database,
	)
	databaseURI += "?pool_max_conns=10&pool_max_conn_lifetime=1m&pool_max_conn_idle_time=1m"
	db, err := pgxpool.Connect(context.Background(), databaseURI)
	if err != nil {
		return nil, err
	}

	queries := []string{
		users,
		schemes,
	}
	for _, v := range queries {
		if _, err := db.Exec(context.Background(), v); err != nil {
			return nil, err
		}
	}
	storage := &Database{
		db: db,
		l:  l,
	}
	return storage, nil
}

// CloseDb Метод закрывает соединение с БД
func (s *Database) CloseDb() error {
	s.db.Close()
	return nil
}

func (s *Database) Exec(sql string, arguments ...any) (pgconn.CommandTag, error) {
	return s.db.Exec(context.Background(), sql, arguments...)
}

func (s *Database) QueryRow(sql string, arguments ...any) pgx.Row {
	return s.db.QueryRow(context.Background(), sql, arguments...)
}

func (s *Database) Query(sql string, arguments ...any) (pgx.Rows, error) {
	return s.db.Query(context.Background(), sql, arguments...)
}
