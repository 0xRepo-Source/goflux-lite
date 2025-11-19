package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Token represents an authentication token
type Token struct {
	ID          string    `json:"id"`
	TokenHash   string    `json:"token_hash"`
	User        string    `json:"user"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Revoked     bool      `json:"revoked"`
}

// TokenStore holds all tokens
type TokenStore struct {
	Tokens []Token `json:"tokens"`
}

func main() {
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("goflux-lite-admin version: 0.1.0-lite")
		return
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create":
		createCommand()
	case "list":
		listCommand()
	case "revoke":
		revokeCommand()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`goflux-lite-admin - Token management tool

USAGE:
  goflux-lite-admin <command> [options]

COMMANDS:
  create -user <name> [-permissions <perms>] [-days <days>] [-file <tokens.json>]
  list [-file <tokens.json>]
  revoke <token_id> [-file <tokens.json>]
  help

OPTIONS:
  -user string         Username for the token (required for create)
  -permissions string  Permissions (comma-separated or * for all, default: *)
  -days int           Token validity in days (default: 30)
  -file string        Token file path (default: tokens.json)

EXAMPLES:
  goflux-lite-admin create -user alice -permissions * -days 365
  goflux-lite-admin create -user bob -permissions upload,download -days 90
  goflux-lite-admin list
  goflux-lite-admin revoke tok_abc123

`)
}

func createCommand() {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	user := fs.String("user", "", "username for the token (required)")
	permissions := fs.String("permissions", "*", "permissions (comma-separated or * for all)")
	days := fs.Int("days", 30, "token validity in days")
	file := fs.String("file", "tokens.json", "token file path")
	fs.Parse(os.Args[2:])

	if *user == "" {
		fmt.Println("Error: -user is required")
		fs.Usage()
		os.Exit(1)
	}

	// Load or create token store
	store := loadOrCreateTokenStore(*file)

	// Parse permissions
	var perms []string
	if *permissions == "*" {
		perms = []string{"*"}
	} else {
		perms = strings.Split(*permissions, ",")
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))

	// Generate ID
	idBytes := make([]byte, 6)
	rand.Read(idBytes)
	id := fmt.Sprintf("tok_%x", idBytes)

	// Create token
	newToken := Token{
		ID:          id,
		TokenHash:   tokenHash,
		User:        *user,
		Permissions: perms,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().AddDate(0, 0, *days),
		Revoked:     false,
	}

	// Add to store and save
	store.Tokens = append(store.Tokens, newToken)
	saveTokenStore(store, *file)

	fmt.Println("Token created successfully!")
	fmt.Println()
	fmt.Printf("Token ID:     %s\n", id)
	fmt.Printf("Token:        %s\n", token)
	fmt.Printf("User:         %s\n", *user)
	fmt.Printf("Permissions:  %v\n", perms)
	fmt.Printf("Expires:      %s\n", newToken.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("⚠️  Save this token! It won't be shown again.")
}

func listCommand() {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	file := fs.String("file", "tokens.json", "token file path")
	fs.Parse(os.Args[2:])

	store := loadOrCreateTokenStore(*file)

	if len(store.Tokens) == 0 {
		fmt.Println("No tokens found.")
		return
	}

	// Header
	fmt.Printf("%-16s %-15s %-30s %-10s %-20s\n",
		"ID", "User", "Permissions", "Status", "Expires")
	fmt.Println(strings.Repeat("─", 95))

	// Tokens
	for _, token := range store.Tokens {
		status := "active"
		if token.Revoked {
			status = "revoked"
		} else if time.Now().After(token.ExpiresAt) {
			status = "expired"
		}

		permsStr := strings.Join(token.Permissions, ",")
		if len(permsStr) > 28 {
			permsStr = permsStr[:25] + "..."
		}

		fmt.Printf("%-16s %-15s %-30s %-10s %-20s\n",
			token.ID,
			token.User,
			permsStr,
			status,
			token.ExpiresAt.Format("2006-01-02 15:04"))
	}
}

func revokeCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: goflux-lite-admin revoke <token_id> [-file <tokens.json>]")
		os.Exit(1)
	}

	tokenID := os.Args[2]
	fs := flag.NewFlagSet("revoke", flag.ExitOnError)
	file := fs.String("file", "tokens.json", "token file path")

	// Parse remaining args (skip token_id)
	if len(os.Args) > 3 {
		fs.Parse(os.Args[3:])
	}

	store := loadOrCreateTokenStore(*file)

	// Find and revoke token
	found := false
	for i := range store.Tokens {
		if store.Tokens[i].ID == tokenID {
			store.Tokens[i].Revoked = true
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Token not found: %s\n", tokenID)
		os.Exit(1)
	}

	saveTokenStore(store, *file)
	fmt.Printf("✓ Token %s has been revoked.\n", tokenID)
}

func loadOrCreateTokenStore(filename string) *TokenStore {
	store := &TokenStore{Tokens: []Token{}}

	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading token file: %v\n", err)
			os.Exit(1)
		}

		if err := json.Unmarshal(data, store); err != nil {
			fmt.Printf("Error parsing token file: %v\n", err)
			os.Exit(1)
		}
	}

	return store
}

func saveTokenStore(store *TokenStore, filename string) {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding tokens: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("Error saving token file: %v\n", err)
		os.Exit(1)
	}
}
