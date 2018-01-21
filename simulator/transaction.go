package simulator

import (
	"time"

	"github.com/satori/go.uuid"
)

type Transaction struct {
	Id        string `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

func NewTransaction(k, v string) *Transaction {
	return &Transaction{
		Id:        uuid.Must(uuid.NewV4()).String(),
		Key:       k,
		Value:     v,
		Timestamp: time.Now().UTC().Unix(),
	}
}
