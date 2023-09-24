package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type PostgresEvent struct {
	db *sqlx.DB
}

func NewEventStorage(db *sqlx.DB) *PostgresEvent {
	return &PostgresEvent{db: db}
}

func (db *PostgresEvent) AddEvent(ctx context.Context) {

}
func (db *PostgresEvent) GetEventByID(ctx context.Context) {
	//По времени сортировка
}
func (db *PostgresEvent) GetNewEvents(ctx context.Context) {
	//По времени сортировка
}
func (db *PostgresEvent) GetStatusEventByID(ctx context.Context) {

}
func (db *PostgresEvent) UpdateStatusEventByID(ctx context.Context) error {
	return nil
}
