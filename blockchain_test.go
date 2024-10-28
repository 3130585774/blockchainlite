package blockchainlite

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func removeFileIfExists(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		err := os.Remove(filePath)
		if err != nil {
			log.Fatalf("Failed to remove existing database file: %v", err)
		}
	} else if !os.IsNotExist(err) {
		log.Fatalf("Error checking if file exists: %v", err)
	}
}

func TestInsertBlock(t *testing.T) {
	dbPath := "test_blockchain"

	bc, err := NewBlockchain(dbPath)
	if err != nil {
		t.Fatalf("Failed to create blockchain: %v", err)
	}
	defer func() {
		bc.Close()
		removeFileIfExists(dbPath + ".db")
	}()

	data1 := map[string]interface{}{
		"title": "First Block",
		"value": 123,
	}
	err = bc.AddBlock(data1)
	if err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}

	data2 := []string{"Second Block", "Another piece of data"}
	err = bc.AddBlock(data2)
	if err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}
}

func TestGetBlockHistory(t *testing.T) {
	dbPath := "test_blockchain"

	bc, err := NewBlockchain(dbPath)
	if err != nil {
		t.Fatalf("Failed to create blockchain: %v", err)
	}
	defer func() {
		bc.Close()
		removeFileIfExists(dbPath + ".db")
	}()

	data1 := map[string]interface{}{
		"title": "First Block",
		"value": 123,
	}
	err = bc.AddBlock(data1)
	if err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}

	data2 := []string{"Second Block", "Another piece of data"}
	err = bc.AddBlock(data2)
	if err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}

	// 获取区块历史记录
	blocks, err := bc.GetBlockHistory()
	if err != nil {
		t.Fatalf("Failed to get block history: %v", err)
	}

	if len(blocks) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(blocks))
	}

	// 检查区块内容
	var dataMap1 map[string]interface{}
	err = json.Unmarshal([]byte(blocks[0].Data), &dataMap1)
	if err != nil {
		t.Fatalf("Failed to unmarshal block data: %v", err)
	}

	title, ok1 := dataMap1["title"].(string)
	value, ok2 := dataMap1["value"].(float64)

	if !ok1 || !ok2 || title != "First Block" || value != 123.00 {
		t.Errorf("Expected data for First Block, got %v", dataMap1)
	}

	var dataSlice []string
	err = json.Unmarshal([]byte(blocks[1].Data), &dataSlice)
	if err != nil {
		t.Fatalf("Failed to unmarshal block data: %v", err)
	}

	if len(dataSlice) != 2 || dataSlice[0] != "Second Block" {
		t.Errorf("Expected data for Second Block, got %v", dataSlice)
	}

}
func TestGetLatestBlock(t *testing.T) {
	dbPath := "test_blockchain"

	bc, err := NewBlockchain(dbPath)
	if err != nil {
		t.Fatalf("Failed to create blockchain: %v", err)
	}
	defer func() {
		bc.Close()
		removeFileIfExists(dbPath + ".db")
	}()

	// 插入一些新区块
	err = bc.AddBlock(map[string]interface{}{"title": "First Block", "value": 123})
	if err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}

	err = bc.AddBlock(map[string]interface{}{"title": "Second Block", "value": 456})
	if err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}

	// 获取最新的区块
	latestBlock, err := bc.GetLatestBlock()
	if err != nil {
		t.Fatalf("Failed to get latest block: %v", err)
	}

	if latestBlock == nil {
		t.Error("Expected latest block, got nil")
	} else {
		var dataMap map[string]interface{}
		err = json.Unmarshal([]byte(latestBlock.Data), &dataMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal latest block data: %v", err)
		}

		title, ok1 := dataMap["title"].(string)
		value, ok2 := dataMap["value"].(float64)

		if !ok1 || !ok2 || title != "Second Block" || value != 456 {
			t.Errorf("Expected latest block data for Second Block, got %v", dataMap)
		}
	}
}
