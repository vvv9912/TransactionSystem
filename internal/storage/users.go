package storage

import (
	"TransactionSystem/internal/constant"
	"TransactionSystem/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type PostgresUsers struct {
	db *sqlx.DB
	//sync.Mutex
}

func NewUsersStorage(db *sqlx.DB) *PostgresUsers {
	return &PostgresUsers{db: db}
}
func (db *PostgresUsers) AddUser(ctx context.Context, User model.Users) (int, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	//
	var id int
	row := conn.QueryRowxContext(
		ctx,
		"INSERT INTO users (wallet_id, currency_code,actual_balance,frozen_balance) VALUES ($1, $2,$3, $4) RETURNING wallet_id",
		User.WalletID,
		User.CurrencyСode,
		User.ActualBalance,
		User.FrozenBalance,
	)
	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (db *PostgresUsers) CheckId(ctx context.Context, WalletID int) (int, error) {
	//db.Lock()
	//defer db.Unlock()
	//log.Print("сработал таймер")
	//time.Sleep(10 * time.Second)
	//defer log.Print("сработал мютекс,таймер офф")
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var id int
	if err := conn.GetContext(ctx, &id, `SELECT wallet_id FROM users WHERE wallet_id = $1`, WalletID); err != nil {
		return 0, err
	}
	return id, err
}
func (db *PostgresUsers) GetActualBalanceByID(ctx context.Context, WalletID int) (float64, error) {
	//db.Lock()
	//defer db.Unlock()
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var balance float64
	if err := conn.GetContext(ctx, &balance, `SELECT actual_balance FROM users WHERE wallet_id= $1`, WalletID); err != nil {
		return 0, err
	}
	return balance, err
}
func (db *PostgresUsers) GetFrozenBalanceByID(ctx context.Context, WalletID int) (float64, error) {
	//db.Lock()
	//defer db.Unlock()
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var balance float64
	if err := conn.GetContext(ctx, &balance, `SELECT frozen_balance FROM users WHERE wallet_id= $1`, WalletID); err != nil {
		return 0, err
	}
	return balance, err
}
func (db *PostgresUsers) AddActualBalanceById(ctx context.Context, WalletID int, account float64) error {
	//db.Lock()
	//defer db.Unlock()
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE users SET actual_balance = (actual_balance + $1) WHERE wallet_id = $2`, account, WalletID)
	return err
}
func (db *PostgresUsers) AddFrozenBalanceById(ctx context.Context, WalletID int64, account float64) error {
	//db.Lock()
	//defer db.Unlock()
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE users SET frozen_balance = (users.frozen_balance + $1) WHERE wallet_id = $2`, account, WalletID)
	return err
}
func (db *PostgresUsers) WithdrawById(ctx context.Context, WalletID int, account float64) error {
	//db.Lock()
	//defer db.Unlock()
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	var ActualBalance float64
	rowsActual := tx.QueryRowContext(ctx, `SELECT actual_balance FROM users WHERE wallet_id= $1`, WalletID)
	err = rowsActual.Scan(&ActualBalance)
	if err != nil {
		return err
	}
	var FrozenBalance float64
	rowsFrozen := tx.QueryRowContext(ctx, `SELECT frozen_balance FROM users WHERE wallet_id= $1`, WalletID)
	err = rowsFrozen.Scan(&FrozenBalance)
	if err != nil {
		return err
	}
	if account > ActualBalance {
		return errors.New(constant.ErrAccountSmall)
	}
	newActual := ActualBalance - account
	_, err = tx.ExecContext(ctx, "UPDATE users SET actual_balance = $1 WHERE wallet_id = $2", newActual, WalletID)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	newFrozen := FrozenBalance + account
	_, err = tx.ExecContext(ctx, "UPDATE users SET frozen_balance = $1 WHERE wallet_id = $2", newFrozen, WalletID)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

//func (db *PostgresUsers) SetActualBalanceById(ctx context.Context, id int64, account float64) error {
//	//db.Lock()
//	//defer db.Unlock()
//	conn, err := db.db.Connx(ctx)
//	if err != nil {
//		return err
//	}
//	defer conn.Close()
//	_, err = conn.ExecContext(ctx, `UPDATE users SET account = $1 WHERE id = $2`, account, id)
//	return err
//}
