package notifier

import (
	"TransactionSystem/internal/constant"
	"TransactionSystem/internal/kafka"
	"TransactionSystem/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type EventsStorage interface {
	GetNewEvents(ctx context.Context) ([]model.Transactions, error)
	UpdateStatusEventByID(ctx context.Context, numTransaction string, status int) error
}
type Cacher interface {
	NewTranscation(t model.Transactions) error
	UpdateTransaction(key string, status int) error
}
type Notifier struct {
	Cacher
	EventsStorage
	KafkaProduce *kafka.Producer
	Partitions   int
	Timer        time.Duration
}

func NewNotifier(cacher Cacher, eventsStorage EventsStorage, partitions int, topic string, timer time.Duration) *Notifier {
	n := &Notifier{
		Cacher:        cacher,
		EventsStorage: eventsStorage,
		Partitions:    partitions,
		Timer:         timer,
	}
	producer := kafka.NewProducer(topic)
	n.KafkaProduce = producer
	//n.KafkaProduce.Topic = constant.Topic_Events
	return n
}
func (n *Notifier) NotifyPending(ctx context.Context) ([]model.Transactions, error) {
	//тут считываем из бд и отправляем1
	trans, err := n.GetNewEvents(ctx)
	if errors.Is(err, sql.ErrNoRows) || len(trans) == 0 {
		err = nil
		return nil, err
	}
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "notifier",
				"func":    "NotifyPending",
				"method":  "GetNewEvents",
			}).Fatalln(err)
		return nil, err
	}
	//Преобразуем значение, ключ

	return trans, nil
}
func (n *Notifier) SendNotification(ctx context.Context, trans model.Transactions) error {
	value, err := json.Marshal(trans)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "notifier",
				"func":    "SendNotification",
				"method":  "json.Marshal(trans[i])",
			}).Fatalln(err)
		return err
	}
	numPartition := (trans.WalletID % 9) + 1 //Нулевой нет партиции
	key := strconv.Itoa(trans.TypeTransaction)
	err = n.KafkaProduce.Produce(value, []byte(key), int32(numPartition))
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "notifier",
				"func":    "SendNotification",
				"method":  "KafkaProduce.Produce",
			}).Fatalln(err)
		return err
	}
	//updatedb status and kafka todo
	err = n.UpdateStatusEventByID(ctx, trans.NumberTransaction, constant.Status_SendKafka)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "notifier",
				"func":    "SendNotification",
				"method":  "UpdateStatusEventByID",
			}).Fatalln(err)
		return err
	}
	//cache
	err = n.UpdateTransaction(trans.NumberTransaction, constant.Status_SendKafka)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"package": "notifier",
				"func":    "SendNotification",
				"method":  "UpdateTransaction",
			}).Fatalln(err)
		return err
	}

	return nil
}
func (n *Notifier) StartNotifyCron(ctx context.Context) error {
	go func() {
		for {
			trans, err := n.NotifyPending(ctx)
			if err != nil {
				logrus.WithFields(
					logrus.Fields{
						"package": "notifier",
						"func":    "StartNotifyCron",
						"method":  "NotifyPending",
					}).Fatalln(err)
				return
			}
			select {
			case <-ctx.Done():
				n.KafkaProduce.Close()
				return
			default:

			}

			for i := range trans {
				go func(k int) {
					err = n.SendNotification(ctx, trans[k])
					if err != nil {
						logrus.WithFields(
							logrus.Fields{
								"package": "notifier",
								"func":    "StartNotifyCron",
								"method":  "SendNotification",
							}).Fatalln(err)
						return
					}
					return
				}(i)
			}
			time.Sleep(n.Timer)
		}
	}()
	return nil
}
