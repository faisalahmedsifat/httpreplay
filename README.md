# httpreplay

HTTP traffic capture and replay tool written in Go. Runs as a forward proxy that records request/response pairs and replays them against a target server for migration testing or shadow traffic verification.

## Features

- **Capture mode** — Starts a forward proxy on a local port, proxies all requests to an upstream target, and records full request/response pairs (headers, bodies, status codes, latency) to an NDJSON file.
- **Replay mode** — Reads a captured NDJSON file and replays each request against the original or a new target URL, comparing response status codes along the way.

## Project Structure

```
.
├── main.go                  # Entry point, dispatches subcommands
├── cmd/
│   ├── capture.go           # CLI handler for "capture"
│   └── replay.go            # CLI handler for "replay"
└── internal/
    ├── proxy/
    │   └── proxy.go         # Reverse proxy with request/response interception
    ├── store/
    │   └── store.go         # Data model (Record struct) and NDJSON I/O
    └── replay/
        └── replayer.go      # Replay engine
```

## Requirements

- Go 1.22.2+

No external dependencies — uses only the Go standard library.

## Install

```bash
go build -o httpreplay .
```

Or run directly without building:

```bash
go run .
```

## Usage

### Capture

Start a proxy on port 8080, forwarding to an upstream API, and record traffic to a file:

```bash
./httpreplay capture --port=8080 --target=https://jsonplaceholder.typicode.com --output=requests.ndjson
```

Then configure your browser or HTTP client to use `http://localhost:8080` as a proxy. Every request/response pair will be written to `requests.ndjson`.

### Replay

Replay a captured file against the original target:

```bash
./httpreplay replay --file=requests.ndjson
```

Replay against a different target (useful for migration testing):

```bash
./httpreplay replay --file=requests.ndjson --target=https://staging-api.example.com
```

Each replayed request prints a status comparison:

```
[GET] Replayed /todos/1 -> New Status: 200 (Recorded was: 304)
```

### Inspect

```bash
./httpreplay inspect --file=requests.ndjson
```

*(Currently a stub — prints "Inspecting..." only.)*

## Data Format

Captured data is stored as newline-delimited JSON (NDJSON). Each line is a JSON object with the following fields:

| Field | Type | Description |
|---|---|---|
| `timestamp` | RFC 3339 time | UTC timestamp of the response |
| `method` | string | HTTP method |
| `url` | string | Full request URL |
| `req_headers` | map | Request headers |
| `req_body` | base64 | Raw request body |
| `status_code` | int | Response status code |
| `res_headers` | map | Response headers |
| `res_body` | base64 | Raw response body |
| `duration` | nanoseconds | Request-response latency |
