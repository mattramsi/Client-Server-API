package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"client-server-api/internal/server/config"
	"client-server-api/pkg/errors"
	"client-server-api/pkg/models"
)

type SQLiteRepository struct {
	db      *sql.DB
	timeout time.Duration
}

func NewSQLiteRepository(cfg config.DatabaseConfig) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", cfg.DSN)
	if err != nil {
		return nil, errors.ErroDatabase(err)
	}

	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, errors.ErroDatabase(err)
	}

	repo := &SQLiteRepository{
		db:      db,
		timeout: cfg.Timeout,
	}

	if err := repo.migrate(); err != nil {
		db.Close()
		return nil, errors.ErroDatabase(err)
	}

	return repo, nil
}

func (r *SQLiteRepository) Save(ctx context.Context, cotacao *models.Cotacao) error {
	ctxDB, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	insertSQL := `
		INSERT INTO cotacoes (code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctxDB, insertSQL,
		cotacao.Code,
		cotacao.Codein,
		cotacao.Name,
		cotacao.High,
		cotacao.Low,
		cotacao.VarBid,
		cotacao.PctChange,
		cotacao.Bid,
		cotacao.Ask,
		cotacao.Timestamp,
		cotacao.CreateDate,
	)

	if err != nil {
		if err == context.DeadlineExceeded {
			return errors.ErroTimeoutContext("salvar cotação no banco", err)
		}
		return errors.ErroDatabase(err)
	}

	return nil
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id int64) (*models.Cotacao, error) {
	ctxDB, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	querySQL := `
		SELECT id, code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date, created_at
		FROM cotacoes
		WHERE id = ?`

	var cotacao models.Cotacao
	err := r.db.QueryRowContext(ctxDB, querySQL, id).Scan(
		&cotacao.ID,
		&cotacao.Code,
		&cotacao.Codein,
		&cotacao.Name,
		&cotacao.High,
		&cotacao.Low,
		&cotacao.VarBid,
		&cotacao.PctChange,
		&cotacao.Bid,
		&cotacao.Ask,
		&cotacao.Timestamp,
		&cotacao.CreateDate,
		&cotacao.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErroNotFound("cotação")
		}
		if err == context.DeadlineExceeded {
			return nil, errors.ErroTimeoutContext("buscar cotação no banco", err)
		}
		return nil, errors.ErroDatabase(err)
	}

	return &cotacao, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepository) migrate() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT,
			codein TEXT,
			name TEXT,
			high TEXT,
			low TEXT,
			var_bid TEXT,
			pct_change TEXT,
			bid TEXT,
			ask TEXT,
			timestamp TEXT,
			create_date TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := r.db.Exec(createTableSQL)
	if err != nil {
		return errors.ErroDatabase(err)
	}

	return nil
}


