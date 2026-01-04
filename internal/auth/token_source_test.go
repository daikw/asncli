package auth

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/99designs/keyring"
)

type fakeStore struct {
	getValue string
	getErr   error
}

func (s fakeStore) Get(service, user string) (string, error) {
	if s.getErr != nil {
		return "", s.getErr
	}
	return s.getValue, nil
}

func (s fakeStore) Set(service, user, token string) error { return nil }
func (s fakeStore) Delete(service, user string) error     { return nil }

func TestTokenSourcePrefersEnv(t *testing.T) {
	if err := os.Setenv(envToken, "env-token"); err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}
	t.Cleanup(func() { _ = os.Unsetenv(envToken) })

	source := NewTokenSource(TokenSourceOptions{Service: "svc", User: "user", Store: fakeStore{getValue: "stored"}})
	got, err := source.Token(context.Background())
	if err != nil {
		t.Fatalf("Token returned unexpected error: %v", err)
	}
	if want := "env-token"; got != want {
		t.Errorf("Token() = %q, want %q (env should take precedence over store)", got, want)
	}
}

func TestTokenSourceUsesStore(t *testing.T) {
	source := NewTokenSource(TokenSourceOptions{Service: "svc", User: "user", Store: fakeStore{getValue: "stored"}})
	got, err := source.Token(context.Background())
	if err != nil {
		t.Fatalf("Token returned unexpected error: %v", err)
	}
	if want := "stored"; got != want {
		t.Errorf("Token() = %q, want %q", got, want)
	}
}

func TestTokenSourceNotFound(t *testing.T) {
	source := NewTokenSource(TokenSourceOptions{Service: "svc", User: "user", Store: fakeStore{getErr: keyring.ErrKeyNotFound}})
	_, err := source.Token(context.Background())
	if !errors.Is(err, ErrNoToken) {
		t.Errorf("Token error = %v, want ErrNoToken", err)
	}
}
