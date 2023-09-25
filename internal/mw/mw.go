package mw

import (
	"TransactionSystem/internal/constant"
	"TransactionSystem/internal/model"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type EventsStorage interface {
	AddEvent(ctx context.Context, t model.Transactions) (int, error)
}
type Cacher interface {
	NewTranscation(t model.Transactions) error
	GetTransaction(key string) (model.Transactions, bool)
}
type MW struct {
	EventsStorage
	Cacher
}

func NewMW(eventsStorage EventsStorage, cacher Cacher) *MW {
	return &MW{EventsStorage: eventsStorage, Cacher: cacher}
}
func (m *MW) Invoice(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var invoice model.Invoice
		sWalletID := ctx.QueryParam("WalletID")
		if sWalletID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected WalletID")
		}
		walletID, err := strconv.Atoi(sWalletID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected WalletID")
		}
		invoice.WalletID = walletID

		sCurrencyСode := ctx.QueryParam("CurrencyСode")
		if sCurrencyСode == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected CurrencyСode")
		}
		CurrencyСode, err := strconv.Atoi(sCurrencyСode)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected CurrencyСode")
		}
		invoice.CurrencyСode = CurrencyСode

		sAmountMoney := ctx.QueryParam("AmountMoney")
		if sAmountMoney == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected AmountMoney")
		}
		AmountMoney, err := strconv.ParseFloat(sAmountMoney, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Uncorrected AmountMoney")
		}
		invoice.AmountMoney = AmountMoney
		ctxBD := context.TODO()

		byteInvoiec, err := json.Marshal(invoice)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "json marshal invoice"}).Fatalf("%v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Error")
		}

		transaction := model.Transactions{
			NumberTransaction: uuid.UUID{},
			WalletID:          invoice.WalletID,
			CreatedAt:         time.Now(),
			Status:            constant.Status_Create,
			TypeTransaction:   constant.Type_Invoice,
			Data:              string(byteInvoiec),
		}

		err = m.NewTranscation(transaction)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "Add transaction to cache"}).Fatalf("%v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Error")
		}
		_, err = m.AddEvent(ctxBD, transaction)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "AddEvent"}).Fatalf("%v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Error")
		}
		err = next(ctx)
		if err != nil {
			logrus.WithFields(logrus.Fields{"func": "InvoiceHandler"}).Fatalf("%v", err)
			return err
		}
		return nil
	}
}
