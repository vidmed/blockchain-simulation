package simulator

import "github.com/satori/go.uuid"

type block struct {
	PrevHash     string         `json:"prev-block-hash"`
	Hash         string         `json:"block-hash"`
	Transactions []*Transaction `json:"transactions"`
}

func newBlock(prev string) *block {
	return &block{
		PrevHash: prev,
		Hash:     uuid.Must(uuid.NewV4()).String(),
	}
}
func (b *block) next() *block {
	if b == nil {
		return newBlock("")
	}
	// to reduce GC work use the same struct
	b.PrevHash = b.Hash
	b.Hash = uuid.Must(uuid.NewV4()).String()
	b.Transactions = nil
	return b
}
