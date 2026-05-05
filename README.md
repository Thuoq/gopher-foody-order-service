# Gopher Foody Order Service

A robust, transactional microservice for managing customer orders and delivery workflows, built with Golang following Clean Architecture and Domain-Driven Design (DDD) principles.

## 🚀 Features

- **Clean Architecture**: Strict separation of concerns (Domain, Ports, UseCases, Infrastructure, Presentation) for high maintainability.
- **Transactional Order Flow**: Atomic creation of Order aggregate including Order Items and Order History within a single database transaction.
- **Data Integrity (Snapshotting)**: Denormalization of food names and prices at the time of purchase to ensure historical accuracy regardless of future price changes.
- **Order Lifecycle Tracking**: Comprehensive audit trail using Order History to track every status transition (Pending, Confirmed, Cancelled, etc.).
- **Resilient Communication**: Built-in HTTP client for fetching real-time food metadata from the Restaurant Service.
- **Dependency Injection**: Modular component management using `uber-go/dig`.
- **Database Migrations**: Versioned SQL migrations for reliable schema evolution.

## 🏗 Project Structure

```text
.
├── cmd/
│   └── server/             # Application entry point (main.go)
├── internal/
│   ├── application/
│   │   └── usecases/       # Single-purpose business logic
│   │       └── order/      # CreateOrder, GetOrderDetail, ListMyOrders
│   ├── core/
│   │   ├── domain/         # Domain entities & constants (Order, Item, History)
│   │   └── ports/          # Interfaces for repositories and external clients
│   ├── infrastructure/
│   │   ├── database/       # GORM implementation and repositories
│   │   └── http/           # HTTP Clients for external services (Restaurant)
│   └── presentation/
│       └── http/           # Gin handlers, DTOs, and routing
├── pkg/
│   ├── logger/             # Zap logger wrapper
│   └── response/           # Standardized API response & validation utility
├── migrations/             # SQL migration files (Up/Down)
└── .env.example            # Environment variables template
```

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Configuration**: Viper
- **Dependency Injection**: Uber-go Dig
- **Logging**: Zap Logger

## 🚦 Getting Started

### 1. Prerequisites
- Go 1.22+
- PostgreSQL
- Running Restaurant Service (for food metadata enrichment)

### 2. Environment Setup
Copy the example environment file and update your database and service connection details:
```bash
cp .env.example .env
```

### 3. Database Migrations
Apply the migrations to set up the Order schema:
```bash
# Using golang-migrate
migrate -path migrations -database "postgres://user:pass@localhost:5432/foody_order_db?sslmode=disable" up
```

### 4. Running the Service
```bash
go mod tidy
go run cmd/server/main.go
```

## 🔌 Inter-service Communication
This service communicates with the **Restaurant Service** via HTTP to fetch and validate food information during the order placement process. Ensure the `RESTAURANT_SERVICE_URL` is correctly configured in your `.env` file.
