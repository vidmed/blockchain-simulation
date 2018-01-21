package simulator

import "github.com/satori/go.uuid"

type Block struct {
	PrevHash     string         `json:"prev-block-hash"`
	Hash         string         `json:"block-hash"`
	Transactions []*Transaction `json:"transactions"`
}

func NewBlock(prev string) *Block {
	return &Block{
		PrevHash: prev,
		Hash:     uuid.Must(uuid.NewV4()).String(),
	}
}
func (b *Block) Next() *Block {
	return NewBlock(b.Hash)
}
