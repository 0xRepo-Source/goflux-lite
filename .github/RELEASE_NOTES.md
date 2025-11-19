# GoFlux Lite v0.1.0 - First Release 🚀

We're excited to announce the first stable release of **GoFlux Lite** - a lightweight, secure file transfer solution with three focused tools.

## 📦 What's Included

### Core Tools
- **🖥️ gfl-server** - Minimal file server with REST API
- **📁 gfl** - Command-line client for file operations  
- **🔑 gfl-admin** - Token management and user administration

### Key Features
- ✅ **Resumable uploads** - Automatic resume for interrupted transfers
- ✅ **Secure authentication** - Token-based with granular permissions
- ✅ **Cross-platform** - Windows, Linux, macOS support
- ✅ **Zero dependencies** - Single binaries, no installation required
- ✅ **Path protection** - Built-in security against traversal attacks

## 🚀 Quick Start

1. **Download** the binaries for your platform
2. **Start server**: `./gfl-server -port 8080`
3. **Create token**: `./gfl-admin create -user admin -permissions *`
4. **Set environment**: `export GOFLUX_TOKEN_LITE="your-token"`
5. **Upload file**: `./gfl put local.txt remote/local.txt`

## 🔧 Use Cases

- **Development teams** - Share builds and artifacts
- **System administrators** - Backup and file distribution
- **Content creators** - Secure file sharing
- **CI/CD pipelines** - Automated file transfers

## 📋 System Requirements

- **Operating System**: Windows 10+, Linux (any recent), macOS 10.14+
- **Architecture**: x64 (amd64)
- **Memory**: 50MB RAM recommended
- **Storage**: Minimal (binaries are ~5-10MB each)

## 🔐 Security

- **Authentication required** by default (can be disabled for testing)
- **Granular permissions** - control upload/download/list access per user
- **Path traversal protection** - prevents `../` directory escape attacks
- **Token expiration** - configurable token lifetimes with automatic cleanup

## 📖 Documentation

- **[README](README.md)** - Quick start and basic usage
- **[Server Guide](docs/gfl-server.md)** - Deployment and configuration
- **[Client Guide](docs/gfl.md)** - File operations and automation
- **[Admin Guide](docs/gfl-admin.md)** - User and token management
- **[Roadmap](docs/roadmap.md)** - Planned features and development

## 🔄 Resumable Uploads

One of our standout features - interrupted uploads automatically resume where they left off:

```bash
# Start upload (may be interrupted)
./gfl put largefile.iso backups/largefile.iso
# ... network issues ...

# Resume automatically on retry
./gfl put largefile.iso backups/largefile.iso
# Continues from where it left off!
```

## 🌟 What Makes GoFlux Lite Special

- **Focused tools** - Each binary does one thing well
- **No web UI bloat** - Pure API-driven architecture
- **Minimal footprint** - Perfect for servers and automation
- **Production ready** - Battle-tested security and reliability
- **Easy deployment** - Copy binaries and run

## ⚡ Performance

- **Fast startup** - Server starts in under 1 second
- **Low memory** - Typically uses under 50MB RAM
- **Concurrent transfers** - Handle multiple uploads/downloads simultaneously
- **Efficient chunks** - Configurable chunk sizes for optimal performance

## 🤝 Community & Support

- **GitHub Issues** - Bug reports and feature requests
- **Documentation** - Comprehensive guides and examples
- **MIT License** - Free for commercial and personal use

## 🔮 What's Next

Check out our [roadmap](docs/roadmap.md) for planned features including:
- Parallel chunk uploads for even faster transfers
- Multiple storage backends (S3, Azure, GCP)
- Enhanced monitoring and metrics
- Additional authentication methods

## 📥 Download

### GitHub Releases
Download pre-built binaries from our [releases page](https://github.com/0xRepo-Source/goflux-lite/releases/tag/v0.1.0).

### Build from Source
```bash
git clone https://github.com/0xRepo-Source/goflux-lite.git
cd goflux-lite
./build.sh  # or .\build.ps1 on Windows
```

## 🙏 Acknowledgments

Thanks to the Go community and everyone who contributed to making this first release possible!

---

**Full Changelog**: https://github.com/0xRepo-Source/goflux-lite/blob/main/CHANGELOG.md