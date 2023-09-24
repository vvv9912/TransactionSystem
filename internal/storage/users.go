package storage

import (
	"TransactionSystem/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type PostgresUsers struct {
	db *sqlx.DB
	//sync.Mutex
}

func NewUsersStorage(db *sqlx.DB) *PostgresUsers {
	return &PostgresUsers{db: db}
}
func (db *PostgresUsers) AddUser(ctx context.Context, User model.Users) (int64, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	//
	var id int64
	row := conn.QueryRowxContext(
		ctx,
		"INSERT INTO users (WalletID, CurrencyСode,ActualBalance,FrozenBalance) VALUES ($1, $2,$3, $4) RETURNING id",
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
	if err := conn.GetContext(ctx, &id, `SELECT WalletID FROM users WHERE WalletID = $1`, WalletID); err != nil {
		return 0, err
	}
	return id, err
}
func (db *PostgresUsers) GetActualBalanceByID(ctx context.Context, WalletID int64) (float64, error) {
	//db.Lock()
	//defer db.Unlock()
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var balance float64
	if err := conn.GetContext(ctx, &balance, `SELECT ActualBalance FROM users WHERE WalletID= $1`, WalletID); err != nil {
		return 0, err
	}
	return balance, err
}
func (db *PostgresUsers) GetFrozenBalanceByID(ctx context.Context, WalletID int64) (float64, error) {
	//db.Lock()
	//defer db.Unlock()
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var balance float64
	if err := conn.GetContext(ctx, &balance, `SELECT FrozenBalance FROM users WHERE WalletID= $1`, WalletID); err != nil {
		return 0, err
	}
	return balance, err
}
func (db *PostgresUsers) AddActualBalanceById(ctx context.Context, WalletID int64, account float64) error {
	//db.Lock()
	//defer db.Unlock()
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE users SET ActualBalance = (ActualBalance + $1) WHERE WalletID = $2`, account, WalletID)
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
	_, err = conn.ExecContext(ctx, `UPDATE users SET FrozenBalance = (FrozenBalance + $1) WHERE WalletID = $2`, account, WalletID)
	return err
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
