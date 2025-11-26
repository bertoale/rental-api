# GO-RENTAL Backend API

## 1. Project Overview

GO-RENTAL adalah backend API untuk sistem manajemen rental kendaraan. API ini menyediakan fitur otentikasi, manajemen user, customer, kendaraan, dan transaksi sewa kendaraan. Dirancang dengan arsitektur RESTful, keamanan JWT, dan dokumentasi Swagger.

**Teknologi utama:**

- Go (Golang)
- Gin Web Framework
- GORM (ORM)
- MySQL
- JWT Authentication
- Swagger (API Docs)

---

## 2. Features

- User authentication (JWT)
- CRUD User, Customer, Vehicle, Rent
- Role-based access (admin, staff)
- Global error handling & validation
- Middleware (auth, CORS, error handler)
- Swagger API documentation
- Seeder admin user

---

## 3. Tech Stack

- **Backend Framework:** Gin (Go)
- **ORM/Database:** GORM, PostgreSQL
- **JWT/Auth:** github.com/golang-jwt/jwt/v5
- **Other Utilities:**
  - github.com/joho/godotenv (env loader)
  - github.com/swaggo/gin-swagger (API docs)
  - Mailjet (email, opsional)

---

## 4. Project Structure

```
server/
├── cmd/                # Entry point (main.go)
├── docs/               # Swagger docs
├── internal/           # Domain logic
│   ├── user/           # User module (CRUD, auth, seeder)
│   ├── customer/       # Customer module
│   ├── vehicle/        # Vehicle module
│   └── rent/           # Rent/transaction module
├── pkg/                # Shared packages
│   ├── config/         # Config & DB connection
│   ├── middlewares/    # Middleware (auth, error)
│   ├── response/       # Response formatter
│   └── validator/      # Custom validation
├── go.mod, go.sum      # Go modules
└── README.md           # This file
```

**Penjelasan:**

- `cmd/`: Main entry point aplikasi
- `internal/`: Kode utama per domain
- `pkg/`: Utilitas & helper global
- `docs/`: Dokumentasi Swagger

---

## 5. Installation

```bash
# Clone repository
$ git clone <repo-url>
$ cd rental-api

# Install dependencies
$ go mod download
```

---

## 6. Environment Variables

Buat file `.env` di root project. Variabel yang diperlukan:

| Variable           | Keterangan                     |
| ------------------ | ------------------------------ |
| DB_HOST            | Host database                  |
| DB_PORT            | Port database                  |
| DB_USER            | Username database              |
| DB_PASSWORD        | Password database              |
| DB_NAME            | Nama database                  |
| DB_SSLMODE         | SSL mode (disable/require)     |
| JWT_SECRET         | Secret key JWT                 |
| JWT_EXPIRES_IN     | Durasi token JWT (misal: 168h) |
| PORT               | Port aplikasi (default: 5000)  |
| NODE_ENV           | development/production         |
| CORS_ORIGIN        | Origin frontend                |
| MAILJET_API_KEY    | (Opsional) API key Mailjet     |
| MAILJET_API_SECRET | (Opsional) Secret Mailjet      |
| MAILJET_PORT       | (Opsional) SMTP port Mailjet   |
| MAILJET_HOST       | (Opsional) SMTP host Mailjet   |
| MAIL_SENDER_EMAIL  | (Opsional) Email pengirim      |
| MAIL_SENDER_NAME   | (Opsional) Nama pengirim       |

**Contoh .env:**

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=go_rental
DB_SSLMODE=disable
JWT_SECRET=your_super_secret_jwt_key
JWT_EXPIRES_IN=168h
PORT=5000
NODE_ENV=development
CORS_ORIGIN=http://localhost:3000
MAILJET_API_KEY=
MAILJET_API_SECRET=
MAILJET_PORT=587
MAILJET_HOST=in-v3.mailjet.com
MAIL_SENDER_EMAIL=noreply@goevent.com
MAIL_SENDER_NAME=GoEvent App
```

---

## 7. Database Setup

- **Migrasi otomatis:**
  Saat aplikasi dijalankan, migrasi tabel berjalan otomatis.
- **Seeder:**
  Admin user otomatis dibuat jika belum ada (pada file internal/user/seeder.go).

---

## 8. Running the App

**Development:**

```bash
$ go run cmd/main.go
```

**Production:**

```bash
$ go build -o app cmd/main.go
$ ./app
```

---

## 9. API Documentation

- Swagger UI: [http://localhost:5000/swagger/index.html](http://localhost:5000/swagger/index.html)

### Endpoint Utama

#### Auth

- `POST /api/auth/login` — Login user
  - Request:
    ```json
    {
      "username": "admin",
      "password": "11111111"
    }
    ```
  - Response:
    ```json
    {
      "success": true,
      "message": "Login success",
      "data": {
        "token": "<jwt-token>",
        "user": { ... }
      }
    }
    ```

#### User

- `GET /api/user/` — List user
- `POST /api/user/` — Register user
- `GET /api/user/{id}` — Detail user
- `PUT /api/user/{id}` — Update user

#### Customer

- `GET /api/customer/` — List customer
- `POST /api/customer/` — Register customer
- `GET /api/customer/{id}` — Detail customer
- `PUT /api/customer/{id}` — Update customer

#### Vehicle

- `GET /api/vehicle/` — List kendaraan
- `POST /api/vehicle/` — Register kendaraan
- `GET /api/vehicle/{id}` — Detail kendaraan
- `PUT /api/vehicle/{id}` — Update kendaraan
- `DELETE /api/vehicle/{id}` — Hapus kendaraan

#### Rent

- `GET /api/rent/` — List transaksi
- `POST /api/rent/` — Buat transaksi
- `GET /api/rent/{id}` — Detail transaksi
- `PUT /api/rent/{id}` — Update transaksi

**Format Response Sukses:**

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

---

## 10. Authentication

- Sistem login menggunakan JWT (Bearer Token)
- Token dikirim via header `Authorization: Bearer <token>` atau cookie
- Middleware `Authenticate` melindungi route
- Role-based access (admin, staff)

---

## 11. Error Handling

- Format error global:

```json
{
  "success": false,
  "message": "error message"
}
```

- Error 401: Token tidak valid/kadaluarsa
- Error 404: Route/data tidak ditemukan
- Error 400: Validasi gagal

---

## 12. Testing (Opsional)

- Belum tersedia unit test. Tambahkan dengan Go testing framework jika diperlukan.

---

## 13. Deployment

- **Docker:**
  - Buat Dockerfile, build image, dan jalankan container
- **Railway/Render/VPS:**
  - Deploy dengan mengatur env dan build command
- **Tips:**
  - Pastikan variabel .env sudah di-setup di environment production
  - Gunakan reverse proxy (Nginx) untuk production

---

## 14. License

MIT License
