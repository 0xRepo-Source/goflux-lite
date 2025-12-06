package chunk

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		size         int
		expectedSize int
	}{
		{
			name:         "default size",
			size:         0,
			expectedSize: 1024 * 1024,
		},
		{
			name:         "negative size",
			size:         -100,
			expectedSize: 1024 * 1024,
		},
		{
			name:         "custom size",
			size:         512,
			expectedSize: 512,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.size)
			if c.Size != tt.expectedSize {
				t.Errorf("expected size %d, got %d", tt.expectedSize, c.Size)
			}
		})
	}
}

func TestChunker_Split(t *testing.T) {
	c := New(10) // Small chunk size for testing

	data := []byte("Hello, World! This is a test.")
	chunks := c.Split(data)

	// Expected 3 chunks: 10 + 10 + 9 bytes
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}

	// Verify chunk IDs
	for i, chunk := range chunks {
		if chunk.ID != i {
			t.Errorf("chunk %d has wrong ID: %d", i, chunk.ID)
		}
	}

	// Verify first chunk
	if string(chunks[0].Data) != "Hello, Wor" {
		t.Errorf("chunk 0 data mismatch: %s", chunks[0].Data)
	}

	// Verify checksums are valid
	for i, chunk := range chunks {
		hash := sha256.Sum256(chunk.Data)
		expected := hex.EncodeToString(hash[:])
		if chunk.Checksum != expected {
			t.Errorf("chunk %d checksum mismatch", i)
		}
	}
}

func TestChunker_Split_EmptyData(t *testing.T) {
	c := New(100)
	data := []byte("")
	chunks := c.Split(data)

	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for empty data, got %d", len(chunks))
	}
}

func TestChunker_Split_ExactSize(t *testing.T) {
	c := New(10)
	data := []byte("1234567890") // Exactly chunk size
	chunks := c.Split(data)

	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(chunks))
	}

	if string(chunks[0].Data) != "1234567890" {
		t.Errorf("chunk data mismatch: %s", chunks[0].Data)
	}
}

func TestChunker_Reassemble(t *testing.T) {
	c := New(10)

	originalData := []byte("Hello, World! This is a test message for chunking.")
	chunks := c.Split(originalData)

	reassembled, err := c.Reassemble(chunks)
	if err != nil {
		t.Fatalf("Reassemble failed: %v", err)
	}

	if !bytes.Equal(reassembled, originalData) {
		t.Errorf("reassembled data doesn't match original")
	}
}

func TestChunker_Reassemble_EmptyChunks(t *testing.T) {
	c := New(100)
	chunks := []Chunk{}

	reassembled, err := c.Reassemble(chunks)
	if err != nil {
		t.Fatalf("Reassemble failed: %v", err)
	}

	if len(reassembled) != 0 {
		t.Errorf("expected empty result, got %d bytes", len(reassembled))
	}
}

func TestChunker_Reassemble_OutOfOrder(t *testing.T) {
	c := New(10)

	// Create chunks manually with wrong order
	chunk1 := Chunk{
		ID:       0,
		Data:     []byte("first"),
		Checksum: calculateChecksum([]byte("first")),
	}
	chunk2 := Chunk{
		ID:       2, // Wrong ID - should be 1
		Data:     []byte("second"),
		Checksum: calculateChecksum([]byte("second")),
	}

	chunks := []Chunk{chunk1, chunk2}

	_, err := c.Reassemble(chunks)
	if err == nil {
		t.Error("expected error for out-of-order chunks")
	}
}

func TestChunker_Reassemble_CorruptedChecksum(t *testing.T) {
	c := New(10)

	originalData := []byte("Test data for checksum validation")
	chunks := c.Split(originalData)

	// Corrupt the checksum of the first chunk
	chunks[0].Checksum = "0000000000000000000000000000000000000000000000000000000000000000"

	// Change enough characters to make it look like a real hash
	chunks[0].Checksum = "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234"

	_, err := c.Reassemble(chunks)
	if err == nil {
		t.Error("expected error for corrupted checksum")
	}
}

func TestChunker_Reassemble_FallbackChecksum(t *testing.T) {
	c := New(10)

	// Create chunk with fallback-style checksum (mostly zeros)
	chunk := Chunk{
		ID:       0,
		Data:     []byte("test"),
		Checksum: "0000000000000000000000000000000000000000000000000000000000000000",
	}

	chunks := []Chunk{chunk}

	// Should allow fallback checksums with warning
	reassembled, err := c.Reassemble(chunks)
	if err != nil {
		t.Fatalf("Reassemble with fallback checksum failed: %v", err)
	}

	if !bytes.Equal(reassembled, []byte("test")) {
		t.Errorf("reassembled data doesn't match")
	}
}

func TestChunker_RoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		chunkSize int
		data      []byte
	}{
		{
			name:      "small data",
			chunkSize: 100,
			data:      []byte("Small test data"),
		},
		{
			name:      "large data",
			chunkSize: 1024,
			data:      bytes.Repeat([]byte("x"), 10000),
		},
		{
			name:      "binary data",
			chunkSize: 512,
			data:      []byte{0x00, 0xFF, 0xAA, 0x55, 0x12, 0x34, 0x56, 0x78},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.chunkSize)

			chunks := c.Split(tt.data)
			reassembled, err := c.Reassemble(chunks)
			if err != nil {
				t.Fatalf("round trip failed: %v", err)
			}

			if !bytes.Equal(reassembled, tt.data) {
				t.Errorf("round trip data mismatch")
			}
		})
	}
}

// Helper function
func calculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
