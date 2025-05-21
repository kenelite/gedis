# 🚀 Gedis Go SDK

A lightweight Go client for interacting with a **Gedis** server — a Redis-like in-memory data store written in Go.

This SDK enables Go developers to connect to Gedis, execute commands, and manage in-memory data structures like strings, lists, and hashes using a Redis-compatible protocol.

---

## 📦 Features

- ✅ Lightweight TCP-based connection to Gedis
- 🔁 RESP (Redis Serialization Protocol) support
- 📄 Easy-to-use API
- 🧵 Thread-safe connections (one per client)
- 💾 Supports:
    - `SET`, `GET`, `DEL`
    - `LPUSH`, `RPUSH`
    - `HSET`, `HGET`
    - Future extensibility (Pub/Sub, ZSET, Auth, etc.)

---

## 📁 Project Structure

```
gedis-go-sdk/
├── client.go # High-level API (GedisClient)
├── conn.go # TCP connection & command handling
├── resp.go # RESP serialization
├── example/
│ └── main.go # Example usage
├── go.mod
└── README.md

```





---

## 🛠️ Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/gedis-go-sdk.git
cd gedis-go-sdk
go mod tidy
```

🚀 Usage Example
```go


package main

import (
	"fmt"
	"log"
	"gedis-go-sdk/gedis"
)

func main() {
	client, err := gedis.NewClient("localhost:6379")
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer client.Close()

	// Set and Get
	_, _ = client.Set("foo", "bar")
	val, _ := client.Get("foo")
	fmt.Println("GET foo =>", val)

	// List operations
	client.LPush("mylist", "a")
	client.RPush("mylist", "b")

	// Hash operations
	client.HSet("user:1", "name", "Alice")
	name, _ := client.HGet("user:1", "name")
	fmt.Println("HGET user:1 name =>", name)
}
```



🧪 Supported Commands

🔤 String Operations

| Command | Method                   |
|:--------|:-------------------------|
| `SET`   | `client.Set(key, value)` |
| `GET`   | `client.Get(key)`        |
| `DEL`   | `client.Del(key)`        |





📃 List Operations

| Command | Method  |
|:----------|:----------|
| `LPUSH` | `client.LPush(key, val)` |
| `RPUSH` | `client.RPush(key, val)` |


🧾 Hash Operations

| Command  | Method |
|:----------|:----------|
| `HSET`  | `client.HSet(key, field, val)` |
| `HGET`  | `client.HGet(key, field)`      |




🛡️ Roadmap
* [ ] Support Pub/Sub
* [ ] Pipelining support
* [ ] Cluster-aware client
* [ ] Auth support
* [ ] TLS/SSL encryption

🤝 Contributing
Pull requests and GitHub issues are welcome! Please feel free to fork this repo and suggest improvements or fixes.

📄 License
MIT License © [Kenelite]


