## Memcached Server in Go

A simple implementation of a memcached server written in Go. This server supports basic caching operations such as set, get, add, replace, append, and prepend.

#### Features
- Basic Caching Operations: set, get, add, replace, append, prepend
- Concurrent Connections: Handles multiple client connections using Goroutines
- Configurable Port: Specify the port for the server to listen on

##### Prerequisites
- Go 1.20 or higher 
##### Installation
1. Clone the Repository

```bash
git clone <repository-url>
cd <repository-directory>
```

##### Run the Server
Start the server by executing:

```bash

go run cmd/server/main.go
```

By default, the server will listen on port 11211. You can specify a different port using the -port flag:

```bash
go run cmd/server/main.go -port 12345
```
##### Usage
- Set a Value: set <key> <flags> <expiry> <byte-count> [noreply]
- Get a Value: get <key>
- Add a Value: add <key> <flags> <expiry> <byte-count> [noreply]
- Replace a Value: replace <key> <flags> <expiry> <byte-count> [noreply]
- Append to a Value: append <key> <flags> <expiry> <byte-count> [noreply]
- Prepend to a Value: prepend <key> <flags> <expiry> <byte-count> [noreply]
##### Example
To set a value:

```bash
echo -e "set mykey 0 3600 9\r\nvaluehere\r\n" | nc localhost 11211
```
To get the value:

```bash
echo -e "get mykey\r\n" | nc localhost 11211
```