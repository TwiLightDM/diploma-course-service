package errs

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDublicateKey   = errors.New("dublicate key")
)
