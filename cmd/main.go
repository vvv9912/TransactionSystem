package main

import (
	cache2 "TransactionSystem/internal/cache"
	"TransactionSystem/internal/config"
	"TransactionSystem/internal/constant"
	"TransactionSystem/internal/kafka"
	"TransactionSystem/internal/notifier"
	"TransactionSystem/internal/server"
	"TransactionSystem/internal/storage"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "main"}).Fatalf("faild to connetct to database: %v", err)
		return
	}
	defer db.Close()
	var (
		usersStorage = storage.NewUsersStorage(db) //подкл бд
		eventStorage = storage.NewEventStorage(db)
		c            = cache2.NewCache()
		notif        = notifier.NewNotifier(c, eventStorage, 10, constant.Topic_Events, time.Duration(time.Second))
		s            = server.NewServer(eventStorage, c, usersStorage)
		//cache        = cache.New(cache.DefaultExpiration, 0)
		//consumer     = kafka.NewConsumer(cache, usersStorage, transactionStorang)
		//producer     = kafka.NewProducer()
	)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	err = notif.StartNotifyCron(ctx)
	if err != nil {

		return
	}
	cons := kafka.NewConsumer(eventStorage, usersStorage, c, constant.Topic_Events)
	err = cons.ConsumerStart(ctx)
	if err != nil {
		return
	}
	err = s.Start(ctx, config.Get().HTTPServer)
	if err != nil {
		return
	}
}
