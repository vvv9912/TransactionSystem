package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Invoice(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Success")
}
func Withdraw(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Success")
}
func ActualBalance(ctx echo.Context) error {
	actualbalance := ctx.Get("actualbalance").(float64)
	resp := fmt.Sprintf("Актуальный баланс: %f", actualbalance)
	return ctx.String(http.StatusOK, resp)
}
func FrozenBalance(ctx echo.Context) error {
	frozenbalance := ctx.Get("frozenbalance").(float64)
	resp := fmt.Sprintf("Замороженный баланс: %f", frozenbalance)
	return ctx.String(http.StatusOK, resp)
}
