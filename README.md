# Go Image Server

A high-performance image server written in Go that handles image uploads, conversions, and serving with WebP optimization. This server provides a RESTful API for managing images with support for both single and batch uploads.

## Features

- 🚀 High-performance image processing and serving
- 🖼️ Automatic WebP conversion for optimal image delivery
- 📦 Support for single and batch image uploads
- 🔒 Secure file handling with path traversal protection
- 🌐 CORS support for cross-origin requests
- ⚡ Efficient file storage management
- 🐳 Docker support for easy deployment
- 🔄 Graceful shutdown handling

## API Endpoints

### Upload Single Image

```http
POST /images/{entityType}
```

### Upload Multiple Images

```http
POST /images/{entityType}/batch
```

### Get Image

```http
GET /images/{entityType}/{uuid}/{filename}
```

### Delete Image

```http
DELETE /images/{entityType}/{uuid}
```

## Installation

### Prerequisites
- Go 1.22 or higher
- Docker (optional)
- libwebp development files

### Local Setup

1. Clone the repository
```bash
git clone https://github.com/yourusername/image-server.git
cd image-server
```

2. Install dependencies
```bash
go mod download
```

3. Run the server
```bash
go run main.go
```

### Docker Setup

1. Build and run using Docker Compose
```bash
docker compose up -d --build
```

## Configuration

The server can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| ENVIRONMENT | Running environment | production |
| SERVER_PORT | Server port | 8080 |
| MAX_FILE_SIZE | Maximum file size in bytes | 10485760 (10MB) |
| SHUTDOWN_TIMEOUT | Graceful shutdown timeout | 10s |
| STORAGE_PATH | Storage path for production | ./data |
| DEV_STORAGE_PATH | Storage path for development | ./dev-data |

## Architecture

The project follows a clean architecture pattern with the following components:

- **Handler**: HTTP request handling and routing
- **Service**: Business logic and image processing
- **Storage**: File system operations and storage management
- **Config**: Application configuration management

## Security Features

- File size limits
- Path traversal protection
- Secure file handling
- Controlled CORS policies

## Development

### Project Structure

```.
├── config/ # Configuration management
├── handler/ # HTTP handlers and middleware
├── service/ # Business logic layer
├── storage/ # Storage implementation
├── scripts/ # Deployment scripts
├── data/ # Production storage
├── dev-data/ # Development storage
└── main.go # Application entry point
```

### Building from Source

```bash
CGO_ENABLED=1 GOOS=linux go build -o main .
```

## Deployment

The project includes GitHub Actions workflow for automated deployment to a VPS. The deployment process includes:

1. Automated builds
2. Docker container management
3. Zero-downtime deployment
4. Automatic cleanup of old artifacts

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- [chai2010/webp](https://github.com/chai2010/webp) for WebP encoding support
- [caarlos0/env](https://github.com/caarlos0/env) for environment configuration