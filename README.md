# post-gres-go

A **minimal, educational key-value database** written in Go.  
Think of it as a super lightweight, in-memory **"Redis-like"** store with a **write-ahead log (WAL)** for durability.

This project is primarily for learning purposes — it demonstrates:
- Concurrency-safe in-memory storage using custom data structures (`HashMap` & `ShardedMap`)
- TCP server that accepts multiple concurrent clients
- A simple text-based protocol (`SET`, `GET`, `DEL`)
- Write-Ahead Logging (WAL) for crash recovery

---

## Features

✅ **In-Memory Store** — supports single-lock and sharded-map implementations  
✅ **Concurrent Clients** — handles multiple simultaneous TCP connections  
✅ **WAL Persistence** — state is rebuilt on restart  
✅ **Simple Protocol** — human-readable commands via `netcat` or any TCP client  
✅ **Configurable** — choose port, shard count, and WAL path via flags  

---

## Installation

Make sure you have Go 1.20+ installed.

```bash
git clone https://github.com/RushabhMehta2005/post-go-res.git
cd post-gres-go
go build -o post-gres-go .
````

---

## Usage

Run the server:

```bash
./post-gres-go -port 4242 -shards 8 -wal ./wal_files/wal_file
```

### Command-line flags

| Flag      | Default                | Description                                 |
| --------- | ---------------------- | ------------------------------------------- |
| `-port`   | `4242`                 | TCP port to listen on                       |
| `-wal`    | `./wal_files/wal_file` | Path to the write-ahead log file            |
| `-shards` | `8`                    | Number of shards (1 = single `HashMap`)     |

---

## Protocol

Clients communicate with the server over plain TCP using simple newline-terminated commands.

| Command | Syntax          | Response                                 |
| ------- | --------------- | ---------------------------------------- |
| **SET** | `SET key value` | `+OK`                                    |
| **GET** | `GET key`       | `+OK value` or `-GET could not find key` |
| **DEL** | `DEL key`       | `+OK`                                    |

### Example session with `nc`

```bash
# Connect to the server
nc localhost 4242

# Set a value
SET hello world
+OK

# Retrieve it
GET hello
+OK world

# Delete it
DEL hello
+OK

# Try to retrieve again
GET hello
-GET could not find hello in store
```

---

## Internal Design & Implementation Details

`post-gres-go` is intentionally simple, but its internals demonstrate several key ideas in database design.

### 1. In-Memory Store

Two implementations of `InMemStore` are provided:

* **HashMap** — A single `map[string]string` protected by a single `sync.RWMutex`.
  Simple and efficient for low-concurrency scenarios.

* **ShardedMap** — Splits the keyspace across `N` independent maps (shards),
  each with its own lock. This reduces contention when many goroutines
  perform reads/writes on different keys.

You can control which implementation is used with the `-shards` flag:

* `-shards 1` → single `HashMap`
* `-shards >1` → `ShardedMap` with DJB2 hashing to pick the shard

### 2. Concurrency Model

* The server accepts TCP connections using `net.Listener.Accept` in a loop.
* Each client connection is handled in its own goroutine (`handleConnection`).
* Within a connection, commands are processed sequentially (no pipelining).
* The in-memory store uses either a single global lock or per-shard locks
  to ensure safe concurrent reads/writes.

### 3. Write-Ahead Log (WAL)

* Each mutation (`SET` or `DEL`) is serialized to a text format:

  ```
  <opLen><op><keyLen><key><valueLen><value>\n
  ```

  Example: `3SET5hello5world\n`
* The WAL file is opened in append mode, so new writes are always added at the end.
* On startup, the entire WAL file is replayed (line by line) to rebuild the in-memory store.
* `Sync()` is called after each log write to ensure durability.

### 4. Failure Recovery

If the server crashes:

* On restart, the WAL is replayed, restoring the state to just before the crash.
* There is currently no WAL compaction (log will grow indefinitely).
* Adding snapshotting/compaction is a planned improvement.

---

## Roadmap / TODO

* [ ] Write client libraries (Go, Python, JS)
* [ ] Implement new open-addressed hash table for faster lookups
* [ ] WAL compaction
* [ ] Benchmarking & performance tuning
* [ ] Proper error propagation and graceful shutdown
