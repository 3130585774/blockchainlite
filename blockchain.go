package blockchainlite

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	_ "modernc.org/sqlite"
	"strconv"
	"sync"
	"time"
)

type Block struct {
	Index     int
	Timestamp int64
	Data      string
	Hash      string
	PrevHash  string
}

type Blockchain struct {
	db *sql.DB
	mu sync.Mutex // 用于并发处理
}

func NewBlock(index int, data string, prevHash string) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(), // 使用 Unix 时间戳
		Data:      data,
		PrevHash:  prevHash,
	}
	block.Hash = block.calculateHash()
	return block
}

func (b *Block) calculateHash() string {
	record := strconv.Itoa(b.Index) + strconv.FormatInt(b.Timestamp, 10) + b.Data + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func NewBlockchain(bcName string) (*Blockchain, error) {

	//if _, err := os.Stat(bcName + ".db"); err == nil {
	//	return nil, fmt.Errorf("blockchain %s already exists", bcName)
	//} else if !os.IsNotExist(err) {
	//	return nil, err
	//}

	db, err := sql.Open("sqlite", bcName+".db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS blocks (block_index INTEGER PRIMARY KEY, timestamp INTEGER, block_data TEXT, hash TEXT, prev_hash TEXT)")
	if err != nil {
		return nil, err
	}

	return &Blockchain{db: db}, nil
}

func (bc *Blockchain) AddBlock(data interface{}) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var lastBlock Block
	row := bc.db.QueryRow("SELECT block_index, timestamp, block_data, hash, prev_hash FROM blocks ORDER BY block_index DESC LIMIT 1")
	err = row.Scan(&lastBlock.Index, &lastBlock.Timestamp, &lastBlock.Data, &lastBlock.Hash, &lastBlock.PrevHash)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newBlock := NewBlock(lastBlock.Index+1, string(jsonData), lastBlock.Hash)

	// 验证新区块的哈希与前一个区块的哈希
	if newBlock.PrevHash != lastBlock.Hash {
		return errors.New("invalid previous hash")
	}

	_, err = bc.db.Exec("INSERT INTO blocks (block_index, timestamp, block_data, hash, prev_hash) VALUES (?, ?, ?, ?, ?)",
		newBlock.Index, newBlock.Timestamp, newBlock.Data, newBlock.Hash, newBlock.PrevHash)
	return err
}

func (bc *Blockchain) GetLatestBlock() (*Block, error) {
	var block Block
	row := bc.db.QueryRow("SELECT block_index, timestamp, block_data, hash, prev_hash FROM blocks ORDER BY block_index DESC LIMIT 1")
	err := row.Scan(&block.Index, &block.Timestamp, &block.Data, &block.Hash, &block.PrevHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 如果没有找到区块，返回 nil
		}
		return nil, err
	}
	return &block, nil
}

func (bc *Blockchain) GetBlockHistory() ([]*Block, error) {
	rows, err := bc.db.Query("SELECT block_index, timestamp, block_data, hash, prev_hash FROM blocks ORDER BY block_index ASC")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var blocks []*Block
	for rows.Next() {
		var block Block
		err := rows.Scan(&block.Index, &block.Timestamp, &block.Data, &block.Hash, &block.PrevHash)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, &block)
	}
	return blocks, nil
}

func (bc *Blockchain) Close() error {
	return bc.db.Close()
}
