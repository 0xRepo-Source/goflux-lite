# GoFlux Lite - Planned Features

This document outlines potential features and enhancements planned for future versions of GoFlux Lite. Features are organized by priority and complexity.

## High Priority (Next Release)

### Performance & Reliability
- **Parallel chunk uploads** - Upload multiple chunks simultaneously for faster transfers
- **Download resume support** - Resume interrupted downloads like uploads
- **Connection retry logic** - Automatic retry with exponential backoff for network issues
- **Bandwidth limiting** - Optional upload/download speed limits

### Enhanced Authentication
- **Role-based permissions** - More granular permission system beyond upload/download/list
- **Session management** - Track active sessions and force logout

### Monitoring & Logging
- **Structured logging** - JSON log format for better parsing
- **Metrics endpoint** - Prometheus-compatible metrics (/metrics)
- **Health check endpoint** - Simple health status for load balancers
- **Transfer statistics** - Track upload/download volumes and speeds

## Medium Priority (Future Releases)

### User Experience
- **File metadata support** - Preserve timestamps, permissions where possible
- **Compression on-the-fly** - Optional gzip compression for transfers
- **File deduplication** - Skip uploading files that already exist (checksum based)

### Storage Enhancements
- **Multiple storage backends** - S3, Azure Blob, Google Cloud Storage support (Unlikely)
- **Storage quotas** - Per-user or global storage limits
- **File expiration** - Automatic cleanup of old files
- **Storage encryption** - Encrypt files at rest

### Administrative Features
- **User management API** - RESTful API for user/token operations
- **Audit logging** - Track all file operations and admin actions
- **Backup/restore tools** - Easy migration and backup utilities

### API Improvements
- **OpenAPI/Swagger documentation** - Auto-generated API documentation
- **Webhooks** - Notify external systems of file events
- **Bulk operations** - Upload/download multiple files in one request
- **Directory synchronization** - Keep local and remote directories in sync

## Low Priority (Future Consideration)

### Advanced Features
- **File sharing links** - Generate temporary download URLs for sharing
- **File versioning** - Keep multiple versions of uploaded files
- **Content-based routing** - Route files to different storage based on type/size
- **Geographic replication** - Multi-region file storage

### Security Enhancements
- **Two-factor authentication** - TOTP/SMS support for admin operations
- **IP allow/deny lists** - Network-based access control
- **Rate limiting** - Prevent abuse with configurable rate limits
- **Virus scanning** - Optional malware detection for uploads

### Protocol & Transport
- **HTTP/3 support** - QUIC protocol for improved performance
- **gRPC API** - Alternative to REST API for better performance
- **WebSocket streaming** - Real-time file transfer status
- **P2P file sharing** - Direct client-to-client transfers

## Implementation Notes

### Backwards Compatibility
All new features will maintain backwards compatibility with existing:
- Configuration files
- API endpoints
- Client command syntax
- Token formats

### Performance Goals
- Server memory usage should remain under 50MB for typical workloads
- Binary sizes should stay under 20MB per component
- Startup time should remain under 1 second
- Support for 1000+ concurrent connections

### Security Principles
- Default to secure configurations
- Principle of least privilege for all permissions
- Regular security audits and updates
- Clear security documentation for all features

## Contributing

Want to help implement these features? Check our priorities and:

1. **Open an issue** to discuss the feature before starting work
2. **Review the roadmap** to see what's already planned
3. **Start with high-priority items** for faster acceptance
4. **Keep the "lite" philosophy** - simple, focused, minimal dependencies

### Feature Request Process
1. Check this roadmap to see if it's already planned
2. Open a GitHub issue with the "enhancement" label
3. Describe the use case and benefits
4. Consider implementation complexity and maintenance burden
5. Discuss with maintainers before starting work

## Version Planning

### v0.1.1 (Next Minor)
- Parallel chunk uploads
- Download resume support
- Structured logging
- Basic metrics endpoint

### v0.1.2
- Multiple storage backends
- Role-based permissions
- Transfer progress API
- File metadata support

### v1.0.0 (Stable Release)
- All high-priority features implemented
- Comprehensive test coverage
- Production deployment guide
- Performance benchmarks

## Explicitly NOT Planned

These features go against the "lite" philosophy and won't be added:

- **GUI client applications** - Keep it command-line focused
- **Built-in web file browser** - API-only server
- **Complex workflow engines** - Simple file operations only
- **Heavy dependencies** - Maintain Go standard library focus
- **Docker container** - Official Docker image for easy deployment
- **Kubernetes operator** - Native Kubernetes deployment and management
- **Terraform provider** - Infrastructure as code for GoFlux deployments
- **Client libraries** - Official libraries for Python, JavaScript, etc.
- **Kitchen sink features** - Each tool does one thing well

## Community Input

Feature priorities may change based on:
- User feedback and requests
- Security requirements
- Performance needs
- Maintenance considerations

Join the discussion on GitHub Issues to influence the roadmap!