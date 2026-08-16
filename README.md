<p align="center">
  <h1 align="center">💰 goCatat API</h1>
  <p align="center">
    <strong>Personal Finance Management REST API — Built with Go & Fiber</strong>
  </p>
  <p align="center">
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
    <a href="https://gofiber.io/"><img src="https://img.shields.io/badge/Fiber-v2.52-00ACD7?style=flat-square&logo=go&logoColor=white" alt="Fiber"></a>
    <a href="https://www.postgresql.org/"><img src="https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL"></a>
    <a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker"></a>
    <a href="https://jwt.io/"><img src="https://img.shields.io/badge/Auth-JWT-000000?style=flat-square&logo=jsonwebtokens&logoColor=white" alt="JWT"></a>
  </p>
</p>

---

## 📖 About

**goCatat** is a personal finance management API that helps users track income and expenses across multiple wallets. Built with clean architecture principles (Repository → Service → Handler), it provides a secure and performant backend for financial tracking applications.

### ✨ Key Features

- 🔐 **JWT Authentication** — Secure login with HTTP-only cookies
- 💳 **Multi-Wallet Support** — Manage cash & non-cash wallets
- 📊 **Transaction Tracking** — Record income & expenses with categories
- 📈 **Financial Summary** — Get income, expense, and balance reports by date range
- 🔍 **Advanced Filtering** — Filter transactions by type, category, and date with pagination
- 🛡️ **Security First** — Bcrypt hashing, parameterized SQL, rate limiting, CORS protection

---

## 🏗️ Architecture

```
.
├── config/             # Database connection & configuration
├── handler/            # HTTP handlers (controllers)
├── middleware/         # JWT auth middleware
├── model/             # Data models & types
├── repository/        # Database queries (data access layer)
├── service/           # Business logic layer
├── utils/             # Helper utilities (JWT generation)
├── main.go            # Application entry point & routing
├── Dockerfile         # Multi-stage Docker build
└── .env.example       # Environment variables template
```

```
Request → Middleware → Handler → Service → Repository → PostgreSQL
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [PostgreSQL](https://www.postgresql.org/download/) 15+
- [Docker](https://www.docker.com/) (optional, for containerized deployment)

### 1. Clone the Repository

```bash
git clone https://github.com/Fauzyfm/gocatat-be.git
cd gocatat-be
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gocatat

JWT_SECRET=your_secret_key_min_32_characters
```

### 3. Setup Database

Create the database and run the schema:

```bash
createdb gocatat
psql -d gocatat -f schema.sql
```

### 4. Run the Server

```bash
go mod download
go run main.go
```

The server will start at `http://localhost:8080`

---

## 🐳 Docker Deployment

```bash
# Build
docker build -t gocatat-api .

# Run
docker run -d \
  --name gocatat-api \
  -p 8080:8080 \
  --env-file .env \
  gocatat-api
```

### EasyPanel Deployment

1. Connect your GitHub repo in EasyPanel
2. EasyPanel auto-detects the `Dockerfile`
3. Add environment variables in the EasyPanel dashboard
4. Set internal port to `8080`
5. Assign your custom domain
6. Deploy 🚀

---

## ⚙️ Environment Variables

| Variable | Description | Default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | — |
| `DB_NAME` | Database name | `gocatat` |
| `JWT_SECRET` | JWT signing secret (min 32 chars) | — |
| `JWT_EXPIRE_HOURS` | JWT token expiration (hours) | `24` |
| `APP_PORT` | Server port | `8080` |
| `APP_ENV` | Environment (`development` / `production`) | `development` |
| `CORS_ORIGINS` | Allowed origins (comma-separated) | `http://localhost:3000` |
| `CORS_METHODS` | Allowed HTTP methods | `GET,POST,PUT,DELETE,PATCH` |
| `CORS_HEADERS` | Allowed headers | `Content-Type,Authorization` |
| `CORS_CREDENTIALS` | Allow credentials | `true` |
| `COOKIE_SECURE` | Secure cookie flag (set `true` for HTTPS) | `false` |
| `COOKIE_SAMESITE` | SameSite cookie policy | `Strict` |
| `COOKIE_DOMAIN` | Cookie domain (e.g. `.example.com`) | — |
| `RATE_LIMIT_MAX` | Max requests per window | `30` |
| `RATE_LIMIT_EXPIRATION_SECONDS` | Rate limit window (seconds) | `60` |

---

## 📡 API Reference

**Base URL:** `https://api.gocatat.my.id/api/v1`

### Health Check

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Service health & DB status |

### Authentication

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Register new user |
| `POST` | `/api/v1/auth/login` | Login & get JWT cookie |
| `POST` | `/api/v1/auth/logout` | Logout & clear cookie |

### Profile 🔒

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/profile` | Get current user info |

### Balance (Wallets) 🔒

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/balance/` | Create a new wallet |
| `GET` | `/api/v1/balance/` | List all wallets |
| `GET` | `/api/v1/balance/:id` | Get wallet by ID |
| `PUT` | `/api/v1/balance/:id` | Update wallet |
| `DELETE` | `/api/v1/balance/:id` | Delete wallet |

### Transactions 🔒

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/transaction/` | Create transaction |
| `GET` | `/api/v1/transaction/` | List transactions (with filters) |
| `GET` | `/api/v1/transaction/summary` | Get financial summary |
| `GET` | `/api/v1/transaction/:id` | Get transaction by ID |
| `PUT` | `/api/v1/transaction/:id` | Update transaction |
| `DELETE` | `/api/v1/transaction/:id` | Delete transaction |

> 🔒 = Requires authentication (JWT cookie or `Authorization: Bearer <token>`)

### Query Parameters (GET /transaction/)

| Parameter | Type | Example | Description |
|---|---|---|---|
| `type` | string | `cash` / `nonCash` | Filter by wallet type |
| `category` | string | `income` / `expense` | Filter by category |
| `start_date` | string | `2026-01-01` | Filter from date (YYYY-MM-DD) |
| `end_date` | string | `2026-12-31` | Filter to date (YYYY-MM-DD) |
| `page` | int | `1` | Page number |
| `limit` | int | `20` | Items per page (max 100) |

### Request / Response Examples

<details>
<summary><strong>Register</strong></summary>

```bash
curl -X POST https://api.gocatat.my.id/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "secret123",
    "confirm_password": "secret123"
  }'
```

```json
{
  "success": true,
  "messagge": "Register Berhasil",
  "data": {
    "id": 1,
    "username": "john",
    "email": "john@example.com",
    "role": "user",
    "createdAt": "2026-08-16T12:00:00Z",
    "updateAt": "2026-08-16T12:00:00Z"
  }
}
```

</details>

<details>
<summary><strong>Login</strong></summary>

```bash
curl -X POST https://api.gocatat.my.id/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

```json
{
  "success": true,
  "messagge": "Login Berhasil"
}
```

> JWT token is set as HTTP-only cookie `access_token`

</details>

<details>
<summary><strong>Create Transaction</strong></summary>

```bash
curl -X POST https://api.gocatat.my.id/api/v1/transaction/ \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "balanceID": 1,
    "type": "cash",
    "amount": 500000,
    "category": "income",
    "description": "Gaji bulan Agustus"
  }'
```

```json
{
  "success": true,
  "data": {
    "id": 1,
    "userID": 1,
    "balanceID": 1,
    "type": "cash",
    "amount": 500000,
    "category": "income",
    "description": "Gaji bulan Agustus",
    "createdAt": "2026-08-16T12:00:00Z"
  }
}
```

</details>

---

## 🔒 Security

- **Password Hashing** — Bcrypt with cost factor 12
- **JWT Authentication** — HTTP-only, Secure, SameSite cookies
- **SQL Injection Prevention** — Parameterized queries throughout
- **Rate Limiting** — Configurable request throttling
- **CORS Protection** — Whitelist-based origin control
- **Error Masking** — Internal errors never exposed to clients
- **No Secrets in Code** — All credentials via environment variables

---

## 🛣️ Roadmap

- [ ] Analytics & reporting dashboard
- [ ] AI-powered financial insights
- [ ] Frontend web app
- [ ] Mobile app (Flutter)

---

## 📄 License

This project is for personal use and learning purposes.

---

<p align="center">
  <sub>Built with ❤️ using Go & Fiber</sub>
</p>
