package auth

import (
	"errors"
	"testing"

	"github.com/99designs/keyring"
)

func TestKeyringStore_UsesExistingRing(t *testing.T) {
	mem := keyring.NewArrayKeyring(nil)
	s := &KeyringStore{
		ring: mem,
		open: func(service string) (keyring.Keyring, error) {
			t.Fatalf("open(%q) should not be called when ring already set", service)
			return nil, nil
		},
	}

	if err := s.Set("svc", "user", "token"); err != nil {
		t.Fatalf("Set returned unexpected error: %v", err)
	}
	got, err := s.Get("svc", "user")
	if err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}
	if want := "token"; got != want {
		t.Errorf("Get() = %q, want %q", got, want)
	}
	if err := s.Delete("svc", "user"); err != nil {
		t.Fatalf("Delete returned unexpected error: %v", err)
	}
	_, err = s.Get("svc", "user")
	if !errors.Is(err, keyring.ErrKeyNotFound) {
		t.Errorf("Get after delete: error = %v, want ErrKeyNotFound", err)
	}
}

func TestKeyringStore_OpenError(t *testing.T) {
	openErr := errors.New("open keyring")
	s := &KeyringStore{
		open: func(service string) (keyring.Keyring, error) {
			return nil, openErr
		},
	}

	if err := s.Set("svc", "user", "token"); !errors.Is(err, openErr) {
		t.Errorf("Set error = %v, want %v", err, openErr)
	}
	if _, err := s.Get("svc", "user"); !errors.Is(err, openErr) {
		t.Errorf("Get error = %v, want %v", err, openErr)
	}
	if err := s.Delete("svc", "user"); !errors.Is(err, openErr) {
		t.Errorf("Delete error = %v, want %v", err, openErr)
	}
}

func TestNewKeyringStore(t *testing.T) {
	s := NewKeyringStore()
	if s == nil {
		t.Fatal("NewKeyringStore returned nil")
	}
	if s.open == nil {
		t.Error("NewKeyringStore open function is nil")
	}
}

func TestKeyringStore_RingHandleCaching(t *testing.T) {
	callCount := 0
	mem := keyring.NewArrayKeyring(nil)
	s := &KeyringStore{
		open: func(service string) (keyring.Keyring, error) {
			callCount++
			return mem, nil
		},
	}

	// First call should invoke open
	if err := s.Set("svc", "user", "token"); err != nil {
		t.Fatalf("Set returned unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("open called %d times on first Set, want 1", callCount)
	}

	// Second call should reuse cached ring
	if _, err := s.Get("svc", "user"); err != nil {
		t.Fatalf("Get returned unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("open called %d times total, want 1 (should reuse cached ring)", callCount)
	}
}
