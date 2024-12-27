package blockchain

import (
	"bytes"
	"errors"
	"time"
)

var (
	ErrInvalidBlockIndex  = errors.New("invalid block index")
	ErrInvalidPrevHash    = errors.New("invalid previous hash")
	ErrInvalidProofOfWork = errors.New("invalid proof of work")
)

type Block struct {
	Index     int       `json:"index"`
	Timestamp time.Time `json:"timestamp"`
	PrevHash  []byte    `json:"prev_hash"`
	Hash      []byte    `json:"hash"`
	Data      []byte    `json:"data"`
	Nonce     int       `json:"nonce"`
}

func NewBlock(index int, prevHash []byte, data []byte) (*Block, error) {
	block := &Block{
		Index:     index,
		Timestamp: time.Now(),
		PrevHash:  prevHash,
		Data:      data,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block, nil
}

func (b *Block) ValidateBlock(prevBlock *Block) error {
	if b.Index != prevBlock.Index+1 {
		return ErrInvalidBlockIndex
	}

	if !bytes.Equal(b.PrevHash, prevBlock.Hash) {
		return ErrInvalidPrevHash
	}

	pow := NewProofOfWork(b)
	if !pow.Validate() {
		return ErrInvalidProofOfWork
	}

	return nil
}
