package model

import (
	"github.com/google/uuid"
	"time"
)

// Хрнатися в бд эвентов
type Transactions struct {
	//	ID                int64
	NumberTransaction uuid.UUID
	WalletID          int
	CreatedAt         time.Time
	Status            int
	TypeTransaction   int
	Data              string
}

type Users struct {
	WalletID      int
	CurrencyСode  int
	ActualBalance float64
	FrozenBalance float64
}
type Invoice struct {
	WalletID     int
	CurrencyСode int
	AmountMoney  float64
}
type Withdraw struct {
	WalletID     int
	CurrencyСode int
	AmountMoney  float64
	toWalletID   int
}
