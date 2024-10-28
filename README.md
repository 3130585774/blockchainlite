### 描述

**BlockchainLite** 是一个轻量级的区块链实现，旨在简化和教育目的。该系统允许用户创建一个安全且防篡改的区块链，使用 SQLite 数据库存储数据。

---

### 主要特性
- **区块结构**：每个区块包含索引、时间戳、数据负载、哈希值和前一个区块的哈希值，以确保区块链的完整性和不可变性。
- **数据库存储**：区块以 SQLite 数据库的形式存储，便于访问和操作。
- **数据序列化**：用户可以以 JSON 格式添加数据到区块链，灵活适应多种应用场景。
- **并发支持**：实现通过互斥锁确保线程安全，允许多个 `goroutine` 在不风险数据损坏的情况下添加区块。
- **历史记录检索**：用户可以轻松检索区块历史，以访问以前的条目及其相关数据。

### 快速开始

```go
package main

import (
	"github.com/3130585774/blockchainlite"
	"log"
)

func main() {
	bcName := "myBlockchain"
	server, err := blockchainlite.NewServer(bcName)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	addr := ":8080"
	if err := server.Start(addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer func() {
		if err := server.Stop(); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
	}()

	select {}
}
```

#### API 端点

1. **添加区块**
    - **请求方法**：`POST`
    - **路径**：`/blocks`
    - **请求体**：
      ```json
      {
        "data": "Your block data here"
      }
      ```
    - **响应**：
        - 成功时返回：
          ```json
          {
            "code": 201,
            "data": "Block added successfully"
          }
          ```
        - 失败时返回：
          ```json
          {
            "code": 400,
            "error": "Error message"
          }
          ```

2. **获取最新区块**
    - **请求方法**：`GET`
    - **路径**：`/blocks/latest`
    - **响应**：
        - 成功时返回：
          ```json
          {
            "code": 200,
            "data": {
              "index": 1,
              "timestamp": "2024-01-01T00:00:00Z",
              "data": "Your block data",
              "hash": "abc123",
              "previous_hash": "xyz789"
            }
          }
          ```
        - 失败时返回：
          ```json
          {
            "code": 404,
            "error": "No blocks found"
          }
          ```

3. **获取区块历史**
    - **请求方法**：`GET`
    - **路径**：`/blocks/history`
    - **响应**：
        - 成功时返回：
          ```json
          {
            "code": 200,
            "data": [
              {
                "index": 1,
                "timestamp": "2024-01-01T00:00:00Z",
                "data": "Your block data",
                "hash": "abc123",
                "previous_hash": "xyz789"
              },
              {
                "index": 2,
                "timestamp": "2024-01-02T00:00:00Z",
                "data": "Another block data",
                "hash": "def456",
                "previous_hash": "abc123"
              }
            ]
          }
          ```

---

### 示例用法
以下是使用 Python、Java、JavaScript 和 Go 访问 **BlockchainLite** API 的示例代码。

### 1. Python 示例

```python
import requests
import json

# 添加区块
url = 'http://localhost:8080/blocks'
data = {'data': 'This is my first block data'}
response = requests.post(url, json=data)
print(response.json())

# 获取最新区块
response = requests.get('http://localhost:8080/blocks/latest')
print(response.json())

# 获取区块历史
response = requests.get('http://localhost:8080/blocks/history')
print(response.json())
```

### 2. Java 示例

```java
import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;

public class BlockchainLiteExample {

    public static void main(String[] args) throws Exception {
        // 添加区块
        String url = "http://localhost:8080/blocks";
        String jsonInputString = "{\"data\": \"This is my first block data\"}";
        
        HttpURLConnection conn = (HttpURLConnection) new URL(url).openConnection();
        conn.setRequestMethod("POST");
        conn.setRequestProperty("Content-Type", "application/json");
        conn.setDoOutput(true);
        
        try(OutputStream os = conn.getOutputStream()) {
            byte[] input = jsonInputString.getBytes("utf-8");
            os.write(input, 0, input.length);
        }
        
        System.out.println("Response Code: " + conn.getResponseCode());
        BufferedReader br = new BufferedReader(new InputStreamReader(conn.getInputStream(), "utf-8"));
        StringBuilder response = new StringBuilder();
        String responseLine;
        
        while ((responseLine = br.readLine()) != null) {
            response.append(responseLine.trim());
        }
        System.out.println(response.toString());

        // 获取最新区块
        conn = (HttpURLConnection) new URL("http://localhost:8080/blocks/latest").openConnection();
        conn.setRequestMethod("GET");
        System.out.println("Response Code: " + conn.getResponseCode());
        br = new BufferedReader(new InputStreamReader(conn.getInputStream(), "utf-8"));
        response = new StringBuilder();
        
        while ((responseLine = br.readLine()) != null) {
            response.append(responseLine.trim());
        }
        System.out.println(response.toString());

        // 获取区块历史
        conn = (HttpURLConnection) new URL("http://localhost:8080/blocks/history").openConnection();
        conn.setRequestMethod("GET");
        System.out.println("Response Code: " + conn.getResponseCode());
        br = new BufferedReader(new InputStreamReader(conn.getInputStream(), "utf-8"));
        response = new StringBuilder();
        
        while ((responseLine = br.readLine()) != null) {
            response.append(responseLine.trim());
        }
        System.out.println(response.toString());
    }
}
```

### 3. JavaScript 示例（使用 Fetch API）

```javascript
// 添加区块
fetch('http://localhost:8080/blocks', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({ data: 'This is my first block data' }),
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));

// 获取最新区块
fetch('http://localhost:8080/blocks/latest')
    .then(response => response.json())
    .then(data => console.log(data))
    .catch(error => console.error('Error:', error));

// 获取区块历史
fetch('http://localhost:8080/blocks/history')
    .then(response => response.json())
    .then(data => console.log(data))
    .catch(error => console.error('Error:', error));
```

### 4. Go 示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    // 添加区块
    url := "http://localhost:8080/blocks"
    data := map[string]string{"data": "This is my first block data"}
    jsonData, _ := json.Marshal(data)

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))

    // 获取最新区块
    resp, err = http.Get("http://localhost:8080/blocks/latest")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    body, _ = ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))

    // 获取区块历史
    resp, err = http.Get("http://localhost:8080/blocks/history")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    body, _ = ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```