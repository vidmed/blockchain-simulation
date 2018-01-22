package simulator

import (
	"time"

	"github.com/satori/go.uuid"
)

// Transaction type represents a blockchain transaction
type Transaction struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

// NewTransaction creates new transaction with given key and value.
// Id and Timestamp fields will be generated automatically.
func NewTransaction(k, v string) *Transaction {
	return &Transaction{
		ID:        uuid.NewV4().String(),
		Key:       k,
		Value:     v,
		Timestamp: time.Now().UTC().Unix(),
	}
}
