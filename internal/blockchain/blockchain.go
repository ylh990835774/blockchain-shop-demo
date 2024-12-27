package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var dbPath = "./storage/db/blockchain"

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() (*Blockchain, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("open blockchain db: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close blockchain db: %v", err)
		}
	}()

	data, err := db.Get([]byte("blockchain"), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			// 创建创世区块
			genesisBlock := &Block{
				Index:     0,
				Timestamp: time.Now(),
				PrevHash:  []byte{},
				Hash:      []byte{},
				Data:      []byte("Genesis Block"),
			}
			bc := &Blockchain{Blocks: []*Block{genesisBlock}}
			if err := bc.SaveToDatabase(db); err != nil {
				return nil, fmt.Errorf("save genesis block: %w", err)
			}
			return bc, nil
		}
		return nil, fmt.Errorf("get blockchain data: %w", err)
	}

	var bc Blockchain
	if err := json.Unmarshal(data, &bc); err != nil {
		return nil, fmt.Errorf("unmarshal blockchain data: %w", err)
	}
	return &bc, nil
}

func (bc *Blockchain) SaveToDatabase(db *leveldb.DB) error {
	data, err := json.Marshal(bc)
	if err != nil {
		return fmt.Errorf("marshal blockchain data: %w", err)
	}

	if err := db.Put([]byte("blockchain"), data, nil); err != nil {
		return fmt.Errorf("save blockchain data: %w", err)
	}

	return nil
}

func (bc *Blockchain) AddBlock(data []byte) (*Block, error) {
	prevBlock := bc.GetLatestBlock()
	newBlock := &Block{
		Index:     len(bc.Blocks),
		Timestamp: time.Now(),
		PrevHash:  prevBlock.Hash,
		Data:      data,
	}
	newBlock.Hash = calculateHash(newBlock)
	bc.Blocks = append(bc.Blocks, newBlock)

	// 保存到数据库
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("open blockchain db: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close blockchain db: %v", err)
		}
	}()

	if err := bc.SaveToDatabase(db); err != nil {
		return nil, fmt.Errorf("save block: %w", err)
	}

	return newBlock, nil
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) Validate() error {
	for i := 1; i < len(bc.Blocks); i++ {
		prevBlock := bc.Blocks[i-1]
		currBlock := bc.Blocks[i]

		// 检查哈希
		if !bytes.Equal(currBlock.Hash, calculateHash(currBlock)) {
			return fmt.Errorf("block %d has invalid hash", currBlock.Index)
		}

		// 检查前一个区块的哈希
		if !bytes.Equal(currBlock.PrevHash, prevBlock.Hash) {
			return fmt.Errorf("block %d has invalid prevHash", currBlock.Index)
		}
	}

	// 检查创世区块
	genesisBlock := bc.Blocks[0]
	if !bytes.Equal(genesisBlock.PrevHash, []byte{}) {
		return fmt.Errorf("genesis block has invalid prevHash")
	}

	return nil
}

func calculateHash(b *Block) []byte {
	record := strconv.Itoa(b.Index) + b.Timestamp.String() + string(b.PrevHash) + string(b.Data) + strconv.Itoa(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	return h.Sum(nil)
}

func (bc *Blockchain) FindBlockByHash(hash []byte) *Block {
	for _, block := range bc.Blocks {
		if bytes.Equal(block.Hash, hash) {
			return block
		}
	}
	return nil
}
