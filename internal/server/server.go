package server

import (
	"TransactionSystem/internal/handler"
	"TransactionSystem/internal/mw"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	echo *echo.Echo
}

func NewServer(events mw.EventsStorage, cache mw.Cacher, users mw.UsersStorage) *Server {
	s := &Server{}
	s.echo = echo.New()
	midl := mw.NewMW(events, cache, users)
	s.echo.POST("/invoice", handler.Invoice, midl.Invoice)
	s.echo.POST("/withdraw", handler.Withdraw, midl.Withdraw)
	s.echo.GET("/actualbalance", handler.ActualBalance, midl.ActualBalance)
	s.echo.GET("/frozenbalance", handler.FrozenBalance, midl.FrozenBalance)
	return s
}
func (s *Server) Start(ctx context.Context, address string) error {
	err := s.echo.Start(address)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "ServerStart"}).Fatalf("Server star error: %v", err)
		return err
	}
	return nil
}
