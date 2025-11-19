package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
  get <remote> <local>  Download a file
  put <local> <remote>  Upload a file
  ls [path]            List files/directories

EXAMPLES:
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
	data, err := client.Download(args[0])
	if err != nil {
		log.Fatalf("Download failed: %v", err)
	}

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
	info, err := os.Stat(args[0])
	if err != nil {
		log.Fatalf("File not found: %v", err)
	}

	// Read file data
	data, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fmt.Printf("Uploading %s (%d bytes)...\n", args[0], info.Size())

	// For simplicity, upload as single chunk
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

	fmt.Printf("✓ Upload complete: %s → %s (%d bytes)\n", args[0], args[1], info.Size())
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
