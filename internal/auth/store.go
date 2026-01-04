package auth

import (
	"fmt"
	"runtime"

	"github.com/99designs/keyring"
)

const (
	DefaultService = "asncli"
	DefaultUser    = "token"
)

type TokenStore interface {
	Get(service, user string) (string, error)
	Set(service, user, token string) error
	Delete(service, user string) error
}

type KeyringStore struct {
	ring keyring.Keyring
	err  error
	open func(service string) (keyring.Keyring, error)
}

func NewKeyringStore() *KeyringStore {
	return &KeyringStore{open: openKeyring}
}

func (s *KeyringStore) Get(service, user string) (string, error) {
	kr, err := s.ringHandle(service)
	if err != nil {
		return "", err
	}
	item, err := kr.Get(user)
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func (s *KeyringStore) Set(service, user, token string) error {
	kr, err := s.ringHandle(service)
	if err != nil {
		return err
	}
	return kr.Set(keyring.Item{
		Key:  user,
		Data: []byte(token),
	})
}

func (s *KeyringStore) Delete(service, user string) error {
	kr, err := s.ringHandle(service)
	if err != nil {
		return err
	}
	return kr.Remove(user)
}

func (s *KeyringStore) ringHandle(service string) (keyring.Keyring, error) {
	if s.ring != nil || s.err != nil {
		return s.ring, s.err
	}
	open := s.open
	if open == nil {
		open = openKeyring
	}
	s.ring, s.err = open(service)
	return s.ring, s.err
}

func openKeyring(service string) (keyring.Keyring, error) {
	if service == "" {
		service = DefaultService
	}
	cfg := keyring.Config{
		ServiceName:              service,
		KeychainTrustApplication: runtime.GOOS == "darwin",
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.WinCredBackend,
			keyring.SecretServiceBackend,
		},
	}
	ring, err := keyring.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("open keyring: %w", err)
	}
	return ring, nil
}
