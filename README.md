# GO-RENTAL

GO-RENTAL is a vehicle rental management REST API built with Go and Gin framework. It supports user authentication, customer management, vehicle management, and rental transactions.

## Features

- User authentication (JWT)
- Customer CRUD
- Vehicle CRUD
- Rent transaction (create, update, complete, cancel)
- Relational data (user, customer, vehicle)
- Error handling and validation

## API Endpoints

### Auth

- `POST /api/auth/login` — Login and get JWT token

### User

- `POST /api/user/` — Create user
- `GET /api/user/` — List users
- `GET /api/user/:id` — Get user by ID
- `PUT /api/user/:id` — Update user

### Customer

- `POST /api/customer/` — Create customer
- `GET /api/customer/` — List customers
- `GET /api/customer/:id` — Get customer by ID
- `PUT /api/customer/:id` — Update customer

### Vehicle

- `POST /api/vehicle/` — Create vehicle
- `GET /api/vehicle/` — List vehicles
- `GET /api/vehicle/:id` — Get vehicle by ID
- `PUT /api/vehicle/:id` — Update vehicle
- `DELETE /api/vehicle/:id` — Delete vehicle

### Rent

- `POST /api/rent/` — Create rent (vehicle rental)
- `GET /api/rent/` — List rents
- `GET /api/rent/:id` — Get rent by ID
- `PUT /api/rent/:id/` — Update rent (complete/cancel)

## How to Run

1. Clone this repository
2. Setup your MySQL database and configure `config/config.go`
3. Run migration automatically on startup
4. Start the server:
   ```bash
   go run cmd/main.go
   ```
5. Access API at `http://localhost:5000`

## Authentication

All endpoints (except login) require JWT authentication. Pass the token in the `Cookie` header:

```
Cookie: token=YOUR_JWT_TOKEN
```

## Example Request

### Create Rent

```bash
curl --location --request POST 'http://localhost:5000/api/rent/' \
--header 'Content-Type: application/json' \
--header 'Cookie: token=YOUR_JWT_TOKEN' \
--data-raw '{
  "customer_id": 1,
  "vehicle_id": 2,
  "notes": "Rental for business trip"
}'
```

### Complete Rent

```bash
curl --location --request PUT 'http://localhost:5000/api/rent/3' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Cookie: token=YOUR_JWT_TOKEN' \
--data-urlencode 'status=completed'
```

## Project Structure

```
cmd/main.go                # Entry point
internal/
  customer/                # Customer module
  rent/                    # Rent module
  user/                    # User module
  vehicle/                 # Vehicle module
pkg/
  config/                  # Config and database
  middlewares/             # Auth and error middlewares
  response/                # Response helpers
```

## License

MIT
