package cache

import (
	"TransactionSystem/internal/model"
	"errors"
	"sync"
)

type Cache struct {
	c map[string]model.Transactions
	m sync.Mutex
}

func NewCache() *Cache {
	var c map[string]model.Transactions
	return &Cache{c: c}
}
func (c *Cache) GetTransaction(key string) (model.Transactions, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	value, found := c.c[key]
	if !found {
		return model.Transactions{}, found
	}
	return value, found
}
func (c *Cache) NewTranscation(t model.Transactions) error {
	c.m.Lock()
	defer c.m.Unlock()
	_, found := c.c[t.NumberTransaction.String()]
	if found {
		return errors.New("Transaction exists")
	}
	c.c[t.NumberTransaction.String()] = t
	return nil
}
func (c *Cache) DeleteTranscation(key string) error {
	c.m.Lock()
	defer c.m.Unlock()

	delete(c.c, key)
	return nil
}
func (c *Cache) UpdateTransaction(key string, status int) error {
	c.m.Lock()
	defer c.m.Unlock()
	value, found := c.GetTransaction(key)
	if !found {
		return errors.New("Transaction not found")
	}
	value.Status = status
	err := c.DeleteTranscation(key)
	if err != nil {
		return err
	}
	err = c.NewTranscation(value)
	if err != nil {
		return err
	}
	return nil
}
