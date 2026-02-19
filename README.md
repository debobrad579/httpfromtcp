# HTTP From TCP

An implimentation of the HTTP/1.1 protocol written in go, built on raw TCP connections. Includes a fully functional scientific calculator web app as a demo.

## Project Structure

```
.
├── internal/
│   ├── http/        # HTTP/1.1 protocol implementation
├── calculator-app/  # Scientific calculator web app
└── cmd/
    ├── httpserver/  # Simple demo HTTP server
    └── tcplistener/ # Raw TCP request logger
```

## How It Works

The server reads raw bytes from TCP connections and manually parses HTTP/1.1 requests into the request line, headers, and body. Responses are written back over the same connection using a `response.Writer`.

## Running

**Calculator app**:
```bash
go run ./calculator-app/
```

**Demo HTTP server**:
```bash
go run ./cmd/httpserver/
```

**Raw TCP listener**:
```bash
go run ./cmd/tcplistener/
```

## Running Tests

```bash
go test ./...
```

## Calculator App

A scientific calculator served at `http://localhost:8080`. The frontend sends expressions to a `/api` POST endpoint, which evaluates them server-side using [go-exprtk](https://github.com/Pramod-Devireddy/go-exprtk).

**Supported operations:**
- Basic arithmetic: `+`, `-`, `×`, `÷`, `^`
- Trig functions: `sin`, `cos`, `tan` and their inverses/hyperbolics (and the constant `π`)
- DEG/RAD mode toggle
- Other functions: `√`, `log`, `ln`, `exp`
