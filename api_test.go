package blockchainlite

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	// 创建一个新的服务器实例
	server, err := NewServer("test")
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 启动服务器
	go server.Start(":8080")
	defer func() {
		server.Stop()
		// 删除测试数据库文件
		if err := os.Remove("test.db"); err != nil {
			t.Errorf("Failed to remove test db file: %v", err)
		}
	}()

	// 测试添加区块
	t.Run("AddBlock", func(t *testing.T) {
		data := map[string]interface{}{
			"data": "Test block data",
		}
		jsonData, _ := json.Marshal(data)

		req, err := http.NewRequest("POST", "/blocks", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()

		server.AddBlockHandler(rec, req)

		resp := rec.Result()
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		var response Response
		json.NewDecoder(resp.Body).Decode(&response)
		if response.Code != http.StatusCreated {
			t.Errorf("Expected code %d, got %d", http.StatusCreated, response.Code)
		}
	})

	// 测试获取最新区块
	t.Run("GetLatestBlock", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/blocks/latest", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()

		server.GetLatestBlockHandler(rec, req)

		resp := rec.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response Response
		json.NewDecoder(resp.Body).Decode(&response)
		if response.Code != http.StatusOK {
			t.Errorf("Expected code %d, got %d", http.StatusOK, response.Code)
		}
	})

	// 测试获取区块历史
	t.Run("GetBlockHistory", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/blocks/history", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()

		server.GetBlockHistoryHandler(rec, req)

		resp := rec.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response Response
		json.NewDecoder(resp.Body).Decode(&response)
		if response.Code != http.StatusOK {
			t.Errorf("Expected code %d, got %d", http.StatusOK, response.Code)
		}
	})
}
