# 🚀 Gedis: A Redis-like In-Memory Database in Go

**Gedis** is a lightweight, Redis-compatible in-memory key-value store written in Go. It supports core Redis features including strings, lists, sets, hashes, sorted sets, TTLs, AOF persistence, and more — with extensibility for advanced features like authentication, CLI tooling, proxy load balancing, and a Go SDK.

---

## ✨ Features

- 🔑 String, List, Set, Hash, Sorted Set (ZSet) data types
- ⏱ TTL / Expiry support for keys
- 📂 AOF (Append-Only File) persistence
- 🔁 Command handlers mapped dynamically
- 🔐 Authentication with `users.conf`
- ⚙️ Configurable via `server.conf`
- 💾 Redis-style RESP protocol parser
- 🧪 Redis-like CLI (`gedis-cli`)
- 🧩 Go SDK for client integration
- 🔄 Proxy server for load balancing & failover
- 📌 Sticky session support (client IP or key hash)

---

## 📦 Project Structure

```

gedis/
├── cmd/ # gedis and gedis-cli main programs
├── internal/
│ ├── handle/ # Command handlers (SET, GET, EXPIRE, etc.)
│ ├── response/ # RESP parser and serializer
│ ├── storage/ # In-memory data and AOF logic
│ ├── conf/ # Config file loading (users.conf, server.conf)
│ ├── connect/ # Connection pool
│ ├── proxy/ # Proxy server with load balancing
│ └── sdk/ # Go SDK for clients
├── users.conf
├── server.conf
└── README.md
```



---

## 🛠 Installation

1. **Clone the repo:**

```bash
git clone https://github.com/yourname/gedis.git
cd gedis
```

2. **Build the server:**
```bash
go build -o gedis ./cmd/gedis
```

3. **Run Gedis:**
```bash
./gedis
```

4. **Use the CLI::**
```bash
go run ./cmd/gedis-cli
```



⚙️ Configuration
users.conf
```ini
[users]
default=xxx
admin=xxx 
guest=guest123
```


server.conf
```ini
[server]
port=6379
aof_enabled=true
aof_file=appendonly.aof
```

🧪 Supported Commands
🔤 String
- [x] SET key value
- [x] GET key
- [x] DEL key
- [x] INCR key
- [x] DECR key
- [x] INCRBY key N
- [x] DECRBY key N

📝 List
- [x] LPUSH key val
- [x] RPUSH key val
- [x] LPOP key
- [x] RPOP key
- [x] LRANGE key start stop

🧾 Hash
- [x] HSET key field value
- [x] HGET key field
- [x] HDEL key field

📦 Set
- [x] SADD key member
- [x] SMEMBERS key
- [x] SREM key member
- [x] SUNION key1 key2
- [x] SINTER key1 key2
- [x] SDIFF key1 key2

🧮 Sorted Set (ZSet)
- [x] ZADD key score member
- [x] ZRANGE key start stop
- [x] ZREM key member

🕒 Expiry
- [x] EXPIRE key seconds

📑 AOF Persistence
- [x] AOF logs all mutating commands (SET, LPUSH, SADD, etc.)
- [x] Automatically replays them on startup to restore state
- [x] File is defined in server.conf


🧰 CLI Usage
```bash
go run ./cmd/gedis-cli
```


Example:
```bash
gedis> SET foo bar
OK
gedis> GET foo
"bar"

```

🧩 Gedis Go SDK
Provides a native Go API to talk to Gedis servers.

Install & Import:

```bash
import "github.com/yourname/gedis/internal/sdk"
```

Example:

```bash
client, _ := sdk.NewClient("localhost:6379")
client.Set("x", "123")
val, _ := client.Get("x")
fmt.Println(val) // "123"
```


🔁 Proxy Server
Run gedis-proxy to forward requests to one of several Gedis nodes:

Supports:

- [x] Load balancing
- [x] Failover
- [x] Sticky sessions based on IP or key

🧠 Future Roadmap
* [ ]  RDB snapshots
* [ ]  Pub/Sub support
* [ ]  Cluster support
* [ ]  Transactions (MULTI / EXEC)
* [ ]  TLS/SSL

🤝 Contributing
PRs and feature ideas are welcome. Please open issues or submit a pull request with tests and explanations.

📄 License
MIT License © [Kenelite]


