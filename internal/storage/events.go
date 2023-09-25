package storage

import (
	"TransactionSystem/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type PostgresEvent struct {
	db *sqlx.DB
}

func NewEventStorage(db *sqlx.DB) *PostgresEvent {
	return &PostgresEvent{db: db}
}

func (db *PostgresEvent) AddEvent(ctx context.Context, t model.Transactions) (int, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	//
	var id int
	row := conn.QueryRowxContext(
		ctx,
		"INSERT INTO events (num_transaction, wallet_id, status, type_transaction, data, created_at) VALUES ($1, $2,$3, $4,$5,&6) RETURNING id",
		t.NumberTransaction,
		t.WalletID,
		t.Status,
		t.TypeTransaction,
		t.Data,
		t.CreatedAt,
	)
	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (db *PostgresEvent) GetEventByID(ctx context.Context) {
	//По времени сортировка
}
func (db *PostgresEvent) GetNewEvents(ctx context.Context) ([]model.Transactions, error) {
	//По времени сортировка, по сттусу
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	status := 0
	defer conn.Close()
	var trans []model.Transactions
	if err := conn.GetContext(ctx, &trans, `SELECT (num_transaction, wallet_id, status, type_transaction, data, created_at) FROM events WHERE status= $1`, status); err != nil {
		return nil, err
	}
	return trans, err
}
func (db *PostgresEvent) GetStatusEventByID(ctx context.Context) {

}
func (db *PostgresEvent) UpdateStatusEventByID(ctx context.Context, walletId int, status int) error {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE events SET status = $1 WHERE wallet_id = $2`, status, walletId)
	return err

	return nil
}
