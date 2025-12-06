// Package errors provides custom error types for the goflux-lite application.
// It defines four error categories: AuthError, StorageError, NetworkError, and ValidationError.
// All errors support error wrapping and can be inspected using errors.Is and errors.As.
package errors

import (
	"errors"
	"fmt"
)

// Error categories for goflux-lite

// AuthError represents authentication and authorization failures.
// It includes a type field to distinguish between different auth error conditions.
type AuthError struct {
	Type    AuthErrorType // Specific category of auth error
	Message string        // Human-readable error description
	Err     error         // Underlying cause, if any
}

// AuthErrorType categorizes different authentication and authorization failures.
type AuthErrorType int

const (
	AuthErrorInvalidToken            AuthErrorType = iota // Token format is invalid
	AuthErrorExpiredToken                                 // Token has passed its expiration time
	AuthErrorRevokedToken                                 // Token has been explicitly revoked
	AuthErrorInsufficientPermissions                      // User lacks required permissions
	AuthErrorInvalidCredentials                           // Username or password is incorrect
)

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("auth error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("auth error: %s", e.Message)
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// NewAuthError creates a new authentication error
func NewAuthError(errType AuthErrorType, message string) *AuthError {
	return &AuthError{
		Type:    errType,
		Message: message,
	}
}

// NewAuthErrorWithCause creates a new authentication error with an underlying cause
func NewAuthErrorWithCause(errType AuthErrorType, message string, err error) *AuthError {
	return &AuthError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// StorageError represents file system and storage operation failures.
// It includes the problematic path and a type to identify the failure category.
type StorageError struct {
	Type    StorageErrorType // Specific category of storage error
	Path    string           // File path associated with the error
	Message string           // Human-readable error description
	Err     error            // Underlying cause, if any
}

// StorageErrorType categorizes different storage operation failures.
type StorageErrorType int

const (
	StorageErrorNotFound         StorageErrorType = iota // Requested file or directory not found
	StorageErrorPathTraversal                            // Path attempts to escape storage root
	StorageErrorPermissionDenied                         // Insufficient permissions to access path
	StorageErrorAlreadyExists                            // File or directory already exists
	StorageErrorInvalidPath                              // Path format is invalid
	StorageErrorIO                                       // I/O operation failed
)

func (e *StorageError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("storage error [%s]: %s: %v", e.Path, e.Message, e.Err)
	}
	return fmt.Sprintf("storage error [%s]: %s", e.Path, e.Message)
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

// NewStorageError creates a new storage error
func NewStorageError(errType StorageErrorType, path, message string) *StorageError {
	return &StorageError{
		Type:    errType,
		Path:    path,
		Message: message,
	}
}

// NewStorageErrorWithCause creates a new storage error with an underlying cause
func NewStorageErrorWithCause(errType StorageErrorType, path, message string, err error) *StorageError {
	return &StorageError{
		Type:    errType,
		Path:    path,
		Message: message,
		Err:     err,
	}
}

// NetworkError represents network communication and transport failures.
// It categorizes errors related to connectivity, timeouts, and protocol issues.
type NetworkError struct {
	Type    NetworkErrorType // Specific category of network error
	Message string           // Human-readable error description
	Err     error            // Underlying cause, if any
}

// NetworkErrorType categorizes different network and transport failures.
type NetworkErrorType int

const (
	NetworkErrorConnection        NetworkErrorType = iota // Failed to establish connection
	NetworkErrorTimeout                                   // Operation exceeded time limit
	NetworkErrorInvalidResponse                           // Server response was malformed
	NetworkErrorServerUnavailable                         // Server is unreachable or down
	NetworkErrorBadRequest                                // Request was rejected by server
)

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("network error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("network error: %s", e.Message)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// NewNetworkError creates a new network error
func NewNetworkError(errType NetworkErrorType, message string) *NetworkError {
	return &NetworkError{
		Type:    errType,
		Message: message,
	}
}

// NewNetworkErrorWithCause creates a new network error with an underlying cause
func NewNetworkErrorWithCause(errType NetworkErrorType, message string, err error) *NetworkError {
	return &NetworkError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// ValidationError represents data validation failures.
// It identifies the specific field that failed validation and provides context.
type ValidationError struct {
	Field   string // Name of the field that failed validation
	Message string // Human-readable error description
	Err     error  // Underlying cause, if any
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error [%s]: %s: %v", e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("validation error [%s]: %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// Helper functions to check error types

// IsAuthError checks if an error is or wraps an AuthError.
// It uses errors.As to unwrap error chains.
func IsAuthError(err error) bool {
	var authErr *AuthError
	return errors.As(err, &authErr)
}

// IsStorageError checks if an error is or wraps a StorageError.
// It uses errors.As to unwrap error chains.
func IsStorageError(err error) bool {
	var storageErr *StorageError
	return errors.As(err, &storageErr)
}

// IsNetworkError checks if an error is or wraps a NetworkError.
// It uses errors.As to unwrap error chains.
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr)
}

// IsValidationError checks if an error is or wraps a ValidationError.
// It uses errors.As to unwrap error chains.
func IsValidationError(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr)
}

// GetAuthErrorType extracts the AuthErrorType from an error.
// Returns the type and true if the error is an AuthError, otherwise returns zero value and false.
func GetAuthErrorType(err error) (AuthErrorType, bool) {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type, true
	}
	return 0, false
}

// GetStorageErrorType extracts the StorageErrorType from an error.
// Returns the type and true if the error is a StorageError, otherwise returns zero value and false.
func GetStorageErrorType(err error) (StorageErrorType, bool) {
	var storageErr *StorageError
	if errors.As(err, &storageErr) {
		return storageErr.Type, true
	}
	return 0, false
}
