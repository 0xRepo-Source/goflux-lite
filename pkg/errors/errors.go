package errors

import (
	"errors"
	"fmt"
)

// Error categories for goflux-lite

// AuthError represents authentication/authorization errors
type AuthError struct {
	Type    AuthErrorType
	Message string
	Err     error
}

type AuthErrorType int

const (
	AuthErrorInvalidToken AuthErrorType = iota
	AuthErrorExpiredToken
	AuthErrorRevokedToken
	AuthErrorInsufficientPermissions
	AuthErrorInvalidCredentials
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

// StorageError represents file storage errors
type StorageError struct {
	Type    StorageErrorType
	Path    string
	Message string
	Err     error
}

type StorageErrorType int

const (
	StorageErrorNotFound StorageErrorType = iota
	StorageErrorPathTraversal
	StorageErrorPermissionDenied
	StorageErrorAlreadyExists
	StorageErrorInvalidPath
	StorageErrorIO
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

// NetworkError represents network/transport errors
type NetworkError struct {
	Type    NetworkErrorType
	Message string
	Err     error
}

type NetworkErrorType int

const (
	NetworkErrorConnection NetworkErrorType = iota
	NetworkErrorTimeout
	NetworkErrorInvalidResponse
	NetworkErrorServerUnavailable
	NetworkErrorBadRequest
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

// ValidationError represents data validation errors
type ValidationError struct {
	Field   string
	Message string
	Err     error
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

// IsAuthError checks if an error is an AuthError
func IsAuthError(err error) bool {
	var authErr *AuthError
	return errors.As(err, &authErr)
}

// IsStorageError checks if an error is a StorageError
func IsStorageError(err error) bool {
	var storageErr *StorageError
	return errors.As(err, &storageErr)
}

// IsNetworkError checks if an error is a NetworkError
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr)
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr)
}

// GetAuthErrorType extracts the AuthErrorType from an error
func GetAuthErrorType(err error) (AuthErrorType, bool) {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type, true
	}
	return 0, false
}

// GetStorageErrorType extracts the StorageErrorType from an error
func GetStorageErrorType(err error) (StorageErrorType, bool) {
	var storageErr *StorageError
	if errors.As(err, &storageErr) {
		return storageErr.Type, true
	}
	return 0, false
}
