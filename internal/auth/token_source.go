package auth

import (
	"context"
	"errors"
	"os"

	"github.com/99designs/keyring"
)

var ErrNoToken = errors.New("no stored token, run 'asn auth login'")

const envToken = "ASNCLI_TOKEN"

type TokenSource interface {
	Token(ctx context.Context) (string, error)
}

type TokenSourceOptions struct {
	Service string
	User    string
	Store   TokenStore
}

type tokenSource struct {
	service string
	user    string
	store   TokenStore
}

func NewTokenSource(opts TokenSourceOptions) TokenSource {
	store := opts.Store
	if store == nil {
		store = NewKeyringStore()
	}
	return &tokenSource{
		service: opts.Service,
		user:    opts.User,
		store:   store,
	}
}

func (t *tokenSource) Token(ctx context.Context) (string, error) {
	// Check ASNCLI_TOKEN environment variable
	if token := os.Getenv(envToken); token != "" {
		return token, nil
	}
	// Fall back to stored token
	token, err := t.store.Get(t.service, t.user)
	if err == nil {
		return token, nil
	}
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return "", ErrNoToken
	}
	return "", err
}
