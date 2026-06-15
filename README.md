# URL Shortener

REST API service for creating short links and redirecting users to original URLs.

The project is written in Go and uses the Echo framework. Links are stored in PostgreSQL through `pgx`, and original URLs are encrypted before being saved to the database.

## Project Structure

```text
.
|-- cmd
|   `-- app
|       `-- main.go
|-- internal
|   |-- config
|   |   `-- config.go
|   |-- crypto
|   |   `-- encryptor.go
|   |-- database
|   |   |-- migrations
|   |   |   `-- 001_create_urls.sql
|   |   |-- migrations.go
|   |   `-- postgres.go
|   |-- httpserver
|   |   |-- handlers
|   |   |   `-- url_handler.go
|   |   |-- routes.go
|   |   `-- server.go
|   `-- storage
|       |-- memory.go
|       `-- postgres.go
|-- .env.example
|-- go.mod
|-- go.sum
`-- web
    |-- src
    |   |-- App.tsx
    |   |-- main.tsx
    |   `-- styles.css
    |-- package.json
    `-- vite.config.ts
```

## Requirements

- Go 1.26.3 or newer
- PostgreSQL

## Configuration

The application can be configured with environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `APP_ADDRESS` | `:8080` | HTTP server address |
| `BASE_URL` | `http://localhost:8080` | Base URL used when returning short links |
| `POSTGRES_HOST` | `localhost` | PostgreSQL host |
| `POSTGRES_PORT` | `5432` | PostgreSQL port |
| `POSTGRES_USER` | `postgres` | PostgreSQL username |
| `POSTGRES_PASSWORD` | `sar58yeaf` | PostgreSQL password |
| `POSTGRES_DB` | `password` | PostgreSQL database name |
| `POSTGRES_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `ENCRYPTION_SECRET` | required | Secret used to encrypt and decrypt stored URLs |

Example `.env` values are available in `.env.example`:

```env
APP_ADDRESS=:8080
BASE_URL=http://localhost:8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=sar58yeaf
POSTGRES_DB=password
POSTGRES_SSLMODE=disable
ENCRYPTION_SECRET=change-this-secret-before-production
```

The same `ENCRYPTION_SECRET` must be used after restart. If it changes, previously saved URLs cannot be decrypted.

Before running the app, create a local `.env` file or set the variables in your shell. The tracked `.env` file is intentionally left empty because it should not contain secrets.

## Database

Migrations are embedded into the application and run automatically on startup.

Current migration creates:

- `schema_migrations` for applied migration tracking
- `urls` for short links, with encrypted original URLs stored in `encrypted_url`

SQL injection protection is handled by parameterized `pgx` queries. User input is passed through query parameters like `$1` and `$2`, not string concatenation.

## Run

Start the API:

```bash
go run ./cmd/app
```

The API will be available at:

```text
http://localhost:8080
```

Start the React frontend:

```bash
cd web
npm install
npm run dev
```

The website will be available at:

```text
http://127.0.0.1:5173
```

During development, Vite proxies `/api` and `/health` requests to `http://localhost:8080`.

## Frontend Build

```bash
cd web
npm run build
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
