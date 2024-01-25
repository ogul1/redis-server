# Redis Server implementation in Go

A Redis Server compatible with the official `redis-cli`. Parser for the Redis Serialization Protocol is provided. Currently, a subset of string and list functionalities are implemented.

`TODO: Add AOF to provide some sort of durability.`

## Usage

---
`go run .` to run the server. The server listens for TCP connections on localhost port 6379. Multiple Redis clients can connect to the server and make requests concurrently.

## Supported Functions

---

### String:
* `GET key`
* `SET key value`
* `APPEND key value`
* `STRLEN key`
* `GETRANGE key start end`
* `DECR key`
* `DECRBY key value`

### List:
* `LPUSH key element [element...]`
* `LRANGE key start end`
* `LLEN key`
* `LINDEX key index`
* `LSET key index element`
* `LREM key count element`
* `LPOS key element`
* `RPUSH key element [element...]`
