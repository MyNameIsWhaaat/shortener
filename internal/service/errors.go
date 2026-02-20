package service

import (
	"errors"

	"github.com/MyNameIsWhaaat/shortener/internal/domain"
)

var (
	ErrInvalidURL      = domain.ErrInvalidURL
	ErrURLNotFound     = domain.ErrURLNotFound
	ErrShortCodeExists = domain.ErrShortCodeExists

	ErrEmptyURL         = errors.New("url cannot be empty")
	ErrInvalidShortCode = errors.New("invalid short code format")
	ErrShortCodeTooLong = errors.New("short code too long")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrURLNotFound)
}

func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrShortCodeExists)
}
