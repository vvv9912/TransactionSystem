package server

import (
	"TransactionSystem/internal/handler"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	echo *echo.Echo
}

func NewServer() *Server {
	s := &Server{}
	s.echo = echo.New()
	s.echo.POST("/invoice", handler.Invoice)
	s.echo.POST("/withdraw", handler.Withdraw)
	s.echo.GET("/actualbalance", handler.ActualBalance)
	s.echo.GET("/frozenbalance", handler.FrozenBalance)
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
