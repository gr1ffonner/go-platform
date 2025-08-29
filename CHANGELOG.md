## [Unreleased]

### Added
- Multi-database support (PostgreSQL, MySQL, ClickHouse)
- Dynamic storage selection via environment variables
- Docker Compose configurations for each database backend
- Repository pattern implementation for all storage types
- Graceful shutdown with proper database connection cleanup
- MinIO S3-compatible object storage integration
- NATS message broker integration
- Redis caching layer
- Swagger API documentation generation
- gRPC service definitions with Protobuf
- Database migration system using Goose
- Comprehensive Makefile with storage-specific commands

### Changed
- Refactored application architecture to support multiple databases
- Updated configuration system to handle multiple DSN formats
- Improved error handling and logging across all components
- Enhanced Docker setup with separate compose files per storage type

### Fixed
- ClickHouse authentication issues with admin user
- Database connection closing during graceful shutdown
- Docker build circular dependency issues
- Environment variable validation for storage selection

### Removed
- MongoDB support (replaced with ClickHouse)
- Single docker-compose.yml file (replaced with storage-specific files)
- Hardcoded database connections

## [1.0.2] - 2025-08-11

### Added
- Initial project structure
- PostgreSQL database integration
- Basic HTTP handlers for dog breed API
- Health check endpoints
- Basic logging configuration




