## Memcached Server

A simple implementation of a memcached server written in Go. This server supports basic caching operations such as set, get, add, replace, append, and prepend.

#### Features
- Basic Caching Operations: set, get, add, replace, append, prepend
- Concurrent Connections: Handles multiple client connections 
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

##### Example
To set a value:

```bash
echo -e "set mykey 0 3600 9\r\value\r\n" | nc localhost 11211
```
To get the value:

```bash
echo -e "get mykey\r\n" | nc localhost 11211
```
