package blockchain

import (
	"encoding/hex"
	"fmt"
)

type Service interface {
	RecordTransaction(data []byte) (string, error)
	GetTransaction(txHash string) ([]byte, error)
}

type service struct {
	chain *Blockchain
}

func NewBlockchainService() (Service, error) {
	chain, err := NewBlockchain()
	if err != nil {
		return nil, fmt.Errorf("create blockchain: %w", err)
	}
	return &service{chain: chain}, nil
}

func (s *service) RecordTransaction(data []byte) (string, error) {
	block, err := s.chain.AddBlock(data)
	if err != nil {
		return "", err
	}
	if err := s.ValidateBlockchain(); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", block.Hash), nil
}

func (s *service) GetTransaction(txHash string) ([]byte, error) {
	hash, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction hash: %w", err)
	}

	block := s.chain.FindBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	return block.Data, nil
}

func (s *service) ValidateBlockchain() error {
	return s.chain.Validate()
}
