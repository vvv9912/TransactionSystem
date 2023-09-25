package kafka

import (
	"TransactionSystem/internal/constant"
	"TransactionSystem/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var ( //todo cfg
	bootstrapservers = "localhost"
	groupid          = "myGroup"
	autooffsetreset  = "earliest"
)

type EventsStorager interface {
	//GetNewEvents(ctx context.Context) ([]model.Transactions, error)
	//GetStatusEventByID(ctx context.Context) (int, error)
	//UpdateStatusEventByID(ctx context.Context) error
}
type UsersStorager interface {
	CheckId(ctx context.Context, WalletID int) (int, error)
	AddUser(ctx context.Context, User model.Users) (int, error)
	AddActualBalanceById(ctx context.Context, WalletID int, account float64) error
}

type Cacher interface {
	NewTranscation(t model.Transactions) error
	GetTransaction(key string) (model.Transactions, bool)
}

type Consumer struct {
	EventsStorage EventsStorager
	UsersStorage  UsersStorager
	C             *kafka.Consumer
	Topic         string
}

func NewConsumer(eventstorager EventsStorager, usersStorager UsersStorager, topic string) *Consumer {
	//&

	kfk, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapservers, //todo cfg
		"group.id":          groupid,
		"auto.offset.reset": autooffsetreset,
	})
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "Consumer",
				"func":    "NewConsumer",
				"method":  "NewConsumer",
			}).Fatalln(err)
	}
	return &Consumer{C: kfk, EventsStorage: eventstorager, UsersStorage: usersStorager, Topic: topic}
}

// ch chan *kafka.Message
func (C *Consumer) ConsumerStart(ctx context.Context) error {
	//err := C.C.SubscribeTopics([]string{"Add", "Sub"}, nil) //todo cfg
	err := C.C.Subscribe(C.Topic, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"func": "ConsumerStart"}).Fatalf("Add to topics: %v", err)
		return err
	}
	go func(ctx context.Context) {
		for {
			msg, err := C.C.ReadMessage(time.Millisecond)
			if err == nil {
				go func() {
					//ch <- msg
					err = C.process(ctx, msg)
					if err != nil {
						logrus.WithFields(logrus.Fields{"func": "ConsumerStart"}).Fatalf("process: %v", err)
						return
					}
					return
				}()
				//
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
			select {
			case <-ctx.Done():
				logrus.WithFields(logrus.Fields{"func": "ConsumerStart"}).Fatalf("faild Consumer")
				return
			default:
			}
		}
	}(ctx)
	//
	return nil

}
func (C *Consumer) process(ctx context.Context, msg *kafka.Message) error {

	msgkey := string(msg.Key)
	key, err := strconv.Atoi(msgkey)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "Consumer",
				"func":    "process",
				"method":  "strconv.Atoi(msgkey)",
			}).Fatalln(err)
		return err
	}
	switch key {
	case constant.Type_Invoice:
		err = C.invoice(ctx, msg)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"package": "Consumer",
					"func":    "process",
					"method":  "C.invoice(msg)",
				}).Fatalln(err)
			return err
		}

	case constant.Type_Withdraw:
		err = C.withdraw(ctx, msg)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"package": "Consumer",
					"func":    "process",
					"method":  "C.invoice(msg)",
				}).Fatalln(err)
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Key no found in func process, KEY: %v", key))
	}

	return nil
}

func (C *Consumer) invoice(ctx context.Context, msg *kafka.Message) error {
	var message model.Transactions
	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "Consumer",
				"func":    "invoice",
				"method":  "json.Unmarshal(msg.Value, &message)",
			}).Fatalln(err)
		return err
	}
	var transaction model.Invoice
	switch message.TypeTransaction {
	case constant.Type_Invoice:
		err = json.Unmarshal([]byte(message.Data), &transaction)
		if err != nil {
			return err
		}
	default:
		return errors.New("type not invoice message")
	}
	_, err = C.UsersStorage.CheckId(ctx, message.WalletID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			logrus.WithField("sql: no rows in result set, add users", nil).Warning(err)
			id, err := C.UsersStorage.AddUser(ctx, model.Users{
				WalletID:      message.WalletID,
				CurrencyСode:  transaction.CurrencyСode,
				ActualBalance: 0,
				FrozenBalance: 0,
			})
			if err != nil {
				return err
			}
			logrus.WithField("Add user with id:", id)
		}
	} else {
		logrus.WithFields(logrus.Fields{"package": "Consumer", "func": "ConsumerStart", "method": "CheckId"}).Fatalln(err)
		return err
	}
	err = C.UsersStorage.AddActualBalanceById(ctx, message.WalletID, transaction.AmountMoney)
	if err != nil {
		logrus.WithFields(logrus.Fields{"package": "Consumer", "func": "ConsumerStart", "method": "AddAccountById"}).Fatalln(err)
		return err
	}
	//todo cache
	//todo transaction
	return nil
}
func (C *Consumer) withdraw(ctx context.Context, msg *kafka.Message) error {
	var message model.Transactions
	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "Consumer",
				"func":    "withdraw",
				"method":  "json.Unmarshal(msg.Value, &message)",
			}).Fatalln(err)
		return err
	}
	return nil
}
