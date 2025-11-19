# Changelog

All notable changes to GoFlux Lite will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Auto-discovery system**
  - UDP broadcast service for server announcements
  - `gfl discover` command to find servers on local network
  - `gfl config <server>` command for automatic client configuration
  - Real-time server status display (auth enabled/disabled)

- **Enhanced progress tracking**
  - Visual progress bar for file uploads > 1MB
  - Real-time transfer speed display (B/s, KB/s, MB/s, GB/s)
  - Human-readable file size formatting
  - Chunked upload progress with chunk count display

- **Automatic firewall configuration** 
  - Windows Firewall rule creation for server and discovery ports
  - Administrator privilege detection and friendly messages
  - Automatic rule checking to prevent duplicates

- **Improved user experience**
  - Colored console output for authentication status
  - Professional progress indicators with Unicode characters
  - Smart upload modes (instant for small files, chunked for large)
  - Enhanced error messages and user guidance

## [0.1.0] - 2025-11-19

### Added
- **Core file transfer functionality**
  - `gfl-server` - Lightweight file server with REST API
  - `gfl` - Command-line client for file operations (put/get/ls)
  - `gfl-admin` - Token management tool for authentication

- **Security features**
  - Token-based authentication with granular permissions
  - Challenge-response authentication protocol
  - Path traversal protection
  - Secure token storage and management

- **Resume functionality**
  - Chunked file uploads with automatic resume
  - Session persistence across server restarts
  - Upload status tracking and missing chunk detection
  - Metadata storage for transfer state

- **Authentication system**
  - Role-based permissions (upload/download/list/*) 
  - Token expiration and revocation
  - Environment variable support (`GOFLUX_TOKEN_LITE`)
  - Admin tools for user management

- **Configuration**
  - JSON-based configuration files
  - Flexible storage and metadata directories
  - Optional TLS/HTTPS support
  - Configurable chunk sizes

- **Documentation**
  - Comprehensive README with quick start guide
  - Individual tool documentation (`docs/`)
  - Feature roadmap and development plans
  - Security best practices guide

### Technical Details
- **Language**: Go 1.21+
- **Architecture**: Three separate binaries for focused functionality
- **Storage**: Local filesystem backend
- **API**: REST endpoints for all operations
- **Authentication**: HMAC-based challenge-response + Bearer tokens
- **Resume**: Chunk-based uploads with bitmap tracking

### Security
- Path sanitization prevents directory traversal attacks
- Token authentication required for all operations (configurable)
- Secure token generation using cryptographic randomness
- No sensitive data logged or exposed

### Performance
- Lightweight binaries (~5-10MB each)
- Low memory footprint
- Concurrent connection support
- Efficient chunked transfers

[0.1.0]: https://github.com/0xRepo-Source/goflux-lite/releases/tag/v0.1.0