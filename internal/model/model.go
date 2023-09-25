package model

import (
	"time"
)

// Хрнатися в бд эвентов
type Transactions struct {
	NumberTransaction string    `db:"num_transaction" json:"num_transaction,omitempty"`
	WalletID          int       `db:"wallet_id" json:"wallet_id,omitempty"`
	Status            int       `db:"status" json:"status,omitempty"`
	TypeTransaction   int       `db:"type_transaction" json:"type_transaction,omitempty"`
	Data              string    `db:"data" json:"data,omitempty"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}

type Users struct {
	WalletID      int     `db:"wallet_id" json:"wallet_id"`
	CurrencyСode  int     `db:"currency_code" json:"currency_code"`
	ActualBalance float64 `db:"actual_balance" json:"actual_balance"`
	FrozenBalance float64 `db:"frozen_balance" json:"frozen_balance"`
}
type Invoice struct {
	WalletID     int     `json:"wallet_id"`
	CurrencyСode int     `json:"currency_code"`
	AmountMoney  float64 `json:"amount_money"`
}
type Withdraw struct {
	WalletID     int     `json:"wallet_id,omitempty"`
	CurrencyСode int     `json:"currency_code,omitempty"`
	AmountMoney  float64 `json:"amount_money,omitempty"`
	ToWalletID   int     `json:"to_wallet_id,omitempty"`
}
