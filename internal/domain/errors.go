package domain

import "errors"

var (
    ErrURLNotFound     = errors.New("url not found")
    ErrShortCodeExists = errors.New("short code already exists")
    ErrInvalidURL      = errors.New("invalid url")
)