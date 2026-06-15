# URL Shortener

REST API service for creating short links and redirecting users to original URLs.

The project is written in Go and uses the Echo framework. At the moment links are stored in memory, so all created short URLs are lost after the application restarts.

## Project Structure

```text
.
|-- cmd
|   `-- app
|       `-- main.go
|-- internal
|   |-- config
|   |   `-- config.go
|   |-- httpserver
|   |   |-- handlers
|   |   |   `-- url_handler.go
|   |   |-- routes.go
|   |   `-- server.go
|   `-- storage
|       `-- memory.go
|-- go.mod
`-- go.sum
```

## Requirements

- Go 1.26.3 or newer

## Configuration

The application can be configured with environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `APP_ADDRESS` | `:8080` | HTTP server address |
| `BASE_URL` | `http://localhost:8080` | Base URL used when returning short links |

## Run

```bash
go run ./cmd/app
```

The API will be available at:

```text
http://localhost:8080
```

## API

### Health Check

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

### Create Short URL

```http
POST /api/v1/shorten
Content-Type: application/json
```

Request body:

```json
{
  "url": "https://example.com"
}
```

Response:

```json
{
  "code": "1",
  "short_url": "http://localhost:8080/1",
  "long_url": "https://example.com"
}
```

### Redirect

```http
GET /:code
```

Example:

```text
GET /1
```

If the code exists, the server returns a `302 Found` redirect to the original URL.

## Tests

```bash
go test ./...
```
