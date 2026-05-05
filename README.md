# Gopher Foody Identity Service

Microservice quản lý xác thực và định danh (Identity) cho hệ thống **Gopher Foody**. Service này được xây dựng trên ngôn ngữ **Golang** với kiến trúc **Clean Architecture** và **Domain-Driven Design (DDD)**.

## 🚀 Công nghệ sử dụng

- **Ngôn ngữ**: Golang 1.26
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **RPC Framework**: [gRPC](https://grpc.io/)
- **Cơ sở dữ liệu**: PostgreSQL với [GORM](https://gorm.io/)
- **Dependency Injection**: [uber-go/dig](https://github.com/uber-go/dig)
- **Quản lý cấu hình**: [Viper](https://github.com/spf13/viper)
- **Logging**: [Zap](https://github.com/uber-go/zap)

## 📁 Cấu trúc dự án

Dự án áp dụng chặt chẽ nguyên tắc Clean Architecture, được phân chia thành các package với trách nhiệm riêng biệt:

```text
.
├── api/                   # Định nghĩa gRPC (proto files)
├── cmd/
│   └── server/            # Entry point của ứng dụng, khởi tạo và chạy server (main.go)
├── internal/              # Mã nguồn private của service
│   ├── application/       # Application layer (Use cases)
│   ├── config/            # Nạp và quản lý cấu hình ứng dụng qua biến môi trường
│   ├── domain/            # Domain layer (Entities, Domain Interfaces) - Core business logic
│   ├── infrastructure/    # Infrastructure layer (Database implementation, Repositories)
│   └── presentation/      # Presentation layer (HTTP Handlers, gRPC Services, Router)
├── pkg/                   # Public packages có thể chia sẻ (Logger, định dạng Response...)
├── .env.example           # File biến môi trường mẫu
└── go.mod                 # Quản lý dependencies của Go
```

## ⚙️ Yêu cầu hệ thống

- Go >= 1.26
- PostgreSQL database đang hoạt động
- Công cụ sinh code gRPC (protoc, protoc-gen-go, protoc-gen-go-grpc) nếu bạn cần sửa file proto.

## 🛠 Hướng dẫn cài đặt và chạy thử

1. **Tải các thư viện phụ thuộc**

   ```bash
   go mod download
   go mod tidy
   ```

2. **Cấu hình biến môi trường**

   Tạo một file `.env` từ file mẫu `.env.example` và cập nhật thông tin cấu hình (chuỗi kết nối Database, Port...).

   ```bash
   cp .env.example .env
   ```

3. **Khởi chạy ứng dụng**

   Ứng dụng sử dụng cơ chế Graceful Shutdown và tự động khởi chạy đồng thời **HTTP Server** (phục vụ REST API cho Client/Frontend) và **gRPC Server** (phục vụ gọi nội bộ giữa các microservices).

   ```bash
   go run cmd/server/main.go
   ```

   Mặc định (nếu chưa đổi trong `.env`):
   - **HTTP Server**: Cổng `8080`
   - **gRPC Server**: Cổng `9090`

## 💡 Các tính năng chính hiện tại

- Đăng nhập/Xác thực (SSO Use Case)
- Lấy thông tin hồ sơ người dùng (Get Profile Handler)
- Sẵn sàng tích hợp gRPC communication

## 🗄️ Quản lý Database Migration

Dự án sử dụng [golang-migrate](https://github.com/golang-migrate/migrate) thay vì GORM AutoMigrate để đảm bảo an toàn và tuân thủ chuẩn Production. Các file schema database được đặt ở thư mục root `migrations/`.

### Luồng hoạt động (Up/Down)

- Các file migration luôn đi theo cặp `.up.sql` và `.down.sql` tương ứng với mỗi thay đổi.
- **Tự động chạy khi khởi động**: Nhờ tính năng `embed` của Go, thư mục `migrations` được biên dịch sẵn vào ứng dụng. Khi service khởi động và kết nối thành công tới PostgreSQL, ứng dụng sẽ tự động apply các file `.up.sql` mới nhất (nếu có).
- **Tránh đụng độ**: `golang-migrate` sẽ tự động xử lý cơ chế locking database, giúp đảm bảo khi chạy trên microservice (với nhiều instance cùng bật lên) sẽ không bị lỗi race condition.

### Cách tạo file Migration mới

Khuyến khích cài đặt `golang-migrate` CLI trên máy cá nhân để dễ thao tác:

```bash
# Cài đặt qua Homebrew (trên macOS)
brew install golang-migrate

# Tạo cặp file up/down mới cho tính năng thêm số điện thoại
migrate create -ext sql -dir migrations -seq add_phone_number
```

Lệnh trên sẽ sinh ra 2 file trong thư mục `migrations/`:
- `000002_add_phone_number.up.sql`: Nơi bạn viết lệnh `ALTER TABLE users ADD COLUMN phone VARCHAR(20);`
- `000002_add_phone_number.down.sql`: Nơi bạn viết lệnh `ALTER TABLE users DROP COLUMN phone;`

Sau khi tạo xong file SQL, bạn chỉ cần chạy lại project (`go run cmd/server/main.go`), ứng dụng sẽ tự động cập nhật database cho bạn.
