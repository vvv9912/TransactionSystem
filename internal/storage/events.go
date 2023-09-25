package storage

import (
	"TransactionSystem/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type PostgresEvent struct {
	db *sqlx.DB
}

func NewEventStorage(db *sqlx.DB) *PostgresEvent {
	return &PostgresEvent{db: db}
}

func (db *PostgresEvent) AddEvent(ctx context.Context, t model.Transactions) error {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	//
	//var id int

	rows, err := conn.QueryContext(
		ctx,
		"INSERT INTO events (num_transaction, wallet_id, status, type_transaction, data, created_at) VALUES ($1, $2,$3, $4,$5,$6) ",
		t.NumberTransaction,
		t.WalletID,
		t.Status,
		t.TypeTransaction,
		t.Data,
		t.CreatedAt,
	)
	defer rows.Close()
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
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
	var status int
	status = 0

	defer conn.Close()
	var trans []model.Transactions

	if err := conn.SelectContext(ctx, &trans, `SELECT 
       num_transaction AS num_transaction, 
       wallet_id AS wallet_id, 
       status AS status, 
       type_transaction AS type_transaction, 
       data AS data,
       created_at AS created_at
		FROM events WHERE status= $1`, status); err != nil {
		return nil, err
	}
	//row := conn.QueryRowxContext(ctx,
	//	)

	return trans, err
}

func (db *PostgresEvent) GetStatusEventByID(ctx context.Context) {

}
func (db *PostgresEvent) UpdateStatusEventByID(ctx context.Context, numTransaction string, status int) error {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE events SET status = $1 WHERE num_transaction = $2`, status, numTransaction)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
