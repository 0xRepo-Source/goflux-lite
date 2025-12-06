package errors

import (
	"errors"
	"testing"
)

func TestAuthError(t *testing.T) {
	err := NewAuthError(AuthErrorInvalidToken, "token is invalid")

	if err.Error() != "auth error: token is invalid" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	if err.Type != AuthErrorInvalidToken {
		t.Errorf("unexpected error type: %v", err.Type)
	}
}

func TestAuthErrorWithCause(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewAuthErrorWithCause(AuthErrorExpiredToken, "token expired", cause)

	if err.Unwrap() != cause {
		t.Error("expected to unwrap to cause error")
	}

	if !errors.Is(err, cause) {
		t.Error("expected errors.Is to find cause")
	}
}

func TestIsAuthError(t *testing.T) {
	authErr := NewAuthError(AuthErrorRevokedToken, "revoked")
	regularErr := errors.New("regular error")

	if !IsAuthError(authErr) {
		t.Error("expected IsAuthError to return true for AuthError")
	}

	if IsAuthError(regularErr) {
		t.Error("expected IsAuthError to return false for regular error")
	}
}

func TestGetAuthErrorType(t *testing.T) {
	err := NewAuthError(AuthErrorInsufficientPermissions, "no permission")

	errType, ok := GetAuthErrorType(err)
	if !ok {
		t.Error("expected GetAuthErrorType to succeed")
	}

	if errType != AuthErrorInsufficientPermissions {
		t.Errorf("expected %v, got %v", AuthErrorInsufficientPermissions, errType)
	}

	// Test with non-auth error
	regularErr := errors.New("regular")
	_, ok = GetAuthErrorType(regularErr)
	if ok {
		t.Error("expected GetAuthErrorType to fail for non-auth error")
	}
}

func TestStorageError(t *testing.T) {
	err := NewStorageError(StorageErrorNotFound, "/path/to/file", "file not found")

	expectedMsg := "storage error [/path/to/file]: file not found"
	if err.Error() != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, err.Error())
	}

	if err.Type != StorageErrorNotFound {
		t.Errorf("unexpected error type: %v", err.Type)
	}

	if err.Path != "/path/to/file" {
		t.Errorf("unexpected path: %s", err.Path)
	}
}

func TestStorageErrorWithCause(t *testing.T) {
	cause := errors.New("disk full")
	err := NewStorageErrorWithCause(StorageErrorIO, "/data", "write failed", cause)

	if err.Unwrap() != cause {
		t.Error("expected to unwrap to cause error")
	}
}

func TestIsStorageError(t *testing.T) {
	storageErr := NewStorageError(StorageErrorPathTraversal, "../etc", "traversal attempt")
	regularErr := errors.New("regular error")

	if !IsStorageError(storageErr) {
		t.Error("expected IsStorageError to return true for StorageError")
	}

	if IsStorageError(regularErr) {
		t.Error("expected IsStorageError to return false for regular error")
	}
}

func TestGetStorageErrorType(t *testing.T) {
	err := NewStorageError(StorageErrorAlreadyExists, "/file", "exists")

	errType, ok := GetStorageErrorType(err)
	if !ok {
		t.Error("expected GetStorageErrorType to succeed")
	}

	if errType != StorageErrorAlreadyExists {
		t.Errorf("expected %v, got %v", StorageErrorAlreadyExists, errType)
	}
}

func TestNetworkError(t *testing.T) {
	err := NewNetworkError(NetworkErrorTimeout, "request timed out")

	if err.Error() != "network error: request timed out" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	if err.Type != NetworkErrorTimeout {
		t.Errorf("unexpected error type: %v", err.Type)
	}
}

func TestNetworkErrorWithCause(t *testing.T) {
	cause := errors.New("connection refused")
	err := NewNetworkErrorWithCause(NetworkErrorConnection, "failed to connect", cause)

	if err.Unwrap() != cause {
		t.Error("expected to unwrap to cause error")
	}
}

func TestIsNetworkError(t *testing.T) {
	netErr := NewNetworkError(NetworkErrorServerUnavailable, "server down")
	regularErr := errors.New("regular error")

	if !IsNetworkError(netErr) {
		t.Error("expected IsNetworkError to return true for NetworkError")
	}

	if IsNetworkError(regularErr) {
		t.Error("expected IsNetworkError to return false for regular error")
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("username", "must not be empty")

	expectedMsg := "validation error [username]: must not be empty"
	if err.Error() != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, err.Error())
	}

	if err.Field != "username" {
		t.Errorf("unexpected field: %s", err.Field)
	}
}

func TestIsValidationError(t *testing.T) {
	valErr := NewValidationError("email", "invalid format")
	regularErr := errors.New("regular error")

	if !IsValidationError(valErr) {
		t.Error("expected IsValidationError to return true for ValidationError")
	}

	if IsValidationError(regularErr) {
		t.Error("expected IsValidationError to return false for regular error")
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test error chain: ValidationError -> AuthError -> base error
	baseErr := errors.New("base error")
	authErr := NewAuthErrorWithCause(AuthErrorInvalidToken, "auth failed", baseErr)

	// Should be able to unwrap to base error
	if !errors.Is(authErr, baseErr) {
		t.Error("expected errors.Is to find base error in chain")
	}

	// Should be able to check error type
	var ae *AuthError
	if !errors.As(authErr, &ae) {
		t.Error("expected errors.As to find AuthError")
	}
}
