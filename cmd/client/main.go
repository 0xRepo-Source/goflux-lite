package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/0xRepo-Source/goflux-lite/pkg/config"
	"github.com/0xRepo-Source/goflux-lite/pkg/transport"
)

func main() {
	configFile := flag.String("config", "goflux.json", "path to configuration file")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("gfl version: 0.1.0-lite")
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create HTTP client
	client := transport.NewHTTPClient(cfg.Client.ServerURL)

	// Set authentication token (environment variable takes precedence over config file)
	token := os.Getenv("GOFLUX_TOKEN_LITE")
	if token == "" {
		token = cfg.Client.Token
	}

	if token != "" {
		client.SetAuthToken(token)
	}

	// Execute command
	command := args[0]
	switch command {
	case "discover":
		doDiscover()
	case "config":
		doConfig(args[1:])
	case "get":
		doGet(client, args[1:])
	case "put":
		doPut(client, args[1:])
	case "ls":
		doList(client, args[1:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`GoFlux Lite - Simple file transfer client

USAGE:
  gfl [options] <command> [args...]

OPTIONS:
  -config string    Configuration file (default "goflux.json")
  -version          Show version

COMMANDS:
  discover              Discover GoFlux servers on local network
  config <server>       Configure client for discovered server
  get <remote> <local>  Download a file
  put <local> <remote>  Upload a file
  ls [path]            List files/directories

EXAMPLES:
  gfl discover
  gfl config 192.168.1.100:8080
  gfl put document.pdf files/document.pdf
  gfl get files/document.pdf downloaded.pdf
  gfl ls files/

NOTE: rm and mkdir not available in lite version

`)
}

func loadConfig(configFile string) (*config.Config, error) {
	// Try to find config file
	configPaths := []string{configFile, "config.json", "../goflux.json"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return config.LoadOrCreateConfig(path)
		}
	}

	// Default config if none found
	return &config.Config{
		Client: config.ClientConfig{
			ServerURL: "http://localhost:8080",
			ChunkSize: 1048576,
		},
	}, nil
}

func doGet(client *transport.HTTPClient, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: get <remote_path> <local_path>")
		os.Exit(1)
	}

	fmt.Printf("Downloading %s...\n", args[0])

	// For downloads, we don't have chunking yet, so just show a simple progress indicator
	fmt.Print("Progress: ")

	data, err := client.Download(args[0])
	if err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	// Simple progress animation during download
	fmt.Print("████████████████████████████████████████████████████")
	fmt.Printf("\n")

	if err := os.WriteFile(args[1], data, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("✓ Download complete: %s → %s (%d bytes)\n", args[0], args[1], len(data))
}

func doPut(client *transport.HTTPClient, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: put <local_path> <remote_path>")
		os.Exit(1)
	}

	// Check if file exists
	_, err := os.Stat(args[0])
	if err != nil {
		log.Fatalf("File not found: %v", err)
	}

	// Read file data
	data, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fileSize := len(data)
	chunkSize := 1024 * 1024 // 1MB chunks

	// For small files, upload as single chunk without progress bar
	if fileSize < chunkSize {
		fmt.Printf("Uploading %s (%d bytes)...\n", args[0], fileSize)

		chunkData := transport.ChunkData{
			Path:     args[1],
			ChunkID:  0,
			Data:     data,
			Checksum: "", // Simplified - no checksum for lite version
			Total:    1,
		}

		if err := client.UploadChunk(chunkData); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}

		fmt.Printf("✓ Upload complete: %s → %s (%d bytes)\n", args[0], args[1], fileSize)
		return
	}

	// For larger files, use chunked upload with progress bar
	totalChunks := (fileSize + chunkSize - 1) / chunkSize
	fmt.Printf("Uploading %s (%d bytes) in %d chunks...\n", args[0], fileSize, totalChunks)

	// Create progress bar and speed tracking
	progressWidth := 50
	startTime := time.Now()

	for i := 0; i < totalChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > fileSize {
			end = fileSize
		}

		chunkData := transport.ChunkData{
			Path:     args[1],
			ChunkID:  i,
			Data:     data[start:end],
			Checksum: "", // Simplified - no checksum for lite version
			Total:    totalChunks,
		}

		if err := client.UploadChunk(chunkData); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}

		// Calculate speed and progress
		elapsed := time.Since(startTime).Seconds()
		progress := float64(i+1) / float64(totalChunks)
		filled := int(progress * float64(progressWidth))

		bar := ""
		for j := 0; j < progressWidth; j++ {
			if j < filled {
				bar += "█"
			} else {
				bar += "░"
			}
		}

		percentage := int(progress * 100)
		uploaded := (i + 1) * chunkSize
		if uploaded > fileSize {
			uploaded = fileSize
		}

		// Calculate and format speed
		var speedStr string
		if elapsed > 0 {
			bytesPerSecond := float64(uploaded) / elapsed
			speedStr = formatSpeed(bytesPerSecond)
		} else {
			speedStr = "calculating..."
		}

		fmt.Printf("\r[%s] %d%% (%s) %s", bar, percentage, formatBytes(uploaded)+"/"+formatBytes(fileSize), speedStr)

		if i == totalChunks-1 {
			fmt.Printf("\n")
		}
	}

	fmt.Printf("✓ Upload complete: %s → %s (%d bytes)\n", args[0], args[1], fileSize)
}

func doList(client *transport.HTTPClient, args []string) {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	files, err := client.List(path)
	if err != nil {
		log.Fatalf("List failed: %v", err)
	}

	if len(files) == 0 {
		fmt.Printf("No files in %s\n", path)
		return
	}

	fmt.Printf("Files in %s:\n", path)
	for _, file := range files {
		fmt.Printf("  %s\n", file)
	}
}

func doDiscover() {
	fmt.Println("Discovering GoFlux servers on local network...")

	discovery := transport.NewDiscoveryClient()
	servers, err := discovery.DiscoverServers()
	if err != nil {
		log.Fatalf("Discovery failed: %v", err)
	}

	fmt.Print(discovery.FormatServerList(servers))
}

func doConfig(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: config <server_address>")
		fmt.Println("Example: config 192.168.1.100:8080")
		os.Exit(1)
	}

	serverAddr := args[0]
	fmt.Printf("Configuring client for server: %s\n", serverAddr)

	discovery := transport.NewDiscoveryClient()
	config, err := discovery.GetServerConfig(serverAddr)
	if err != nil {
		log.Fatalf("Failed to get server config: %v", err)
	}

	// Create goflux.json configuration
	clientConfig := map[string]interface{}{
		"client": map[string]interface{}{
			"server_url": fmt.Sprintf("http://%s", serverAddr),
			"chunk_size": 1048576,
			"token":      "", // User must set this manually if auth is required
		},
	}

	// Write configuration to file
	configJSON, err := json.MarshalIndent(clientConfig, "", "  ")
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	if err := os.WriteFile("goflux.json", configJSON, 0644); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	fmt.Println("✓ Configuration saved to goflux.json")

	// Show auth info if required
	if serverConfig, ok := config["server"].(map[string]interface{}); ok {
		if authEnabled, ok := serverConfig["auth_enabled"].(bool); ok && authEnabled {
			fmt.Println()
			fmt.Println("⚠️  This server requires authentication.")
			fmt.Println("   Set GOFLUX_TOKEN_LITE environment variable or edit goflux.json")
			fmt.Println("   Contact the server administrator for a token.")
		}
	}
}

// formatBytes formats byte counts in human-readable format
func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatSpeed formats transfer speed in human-readable format
func formatSpeed(bytesPerSecond float64) string {
	const unit = 1024
	if bytesPerSecond < unit {
		return fmt.Sprintf("%.0f B/s", bytesPerSecond)
	}
	div, exp := float64(unit), 0
	for n := bytesPerSecond / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB/s", bytesPerSecond/div, "KMGTPE"[exp])
}
