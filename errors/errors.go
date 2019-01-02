package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	_                                               = iota
	CodeExistingName                     codes.Code = 10000 + iota // 10000+1
	CodeNonexistentName                                            // +2
	CodeInvalidName                                                // +3
	CodeNotImplemented
	CodeInvalidKey
)

var (
	ErrInvalidName                      = statusError(CodeInvalidName, "invalid account name")
	ErrInvalidKey                       = statusError(CodeInvalidKey, "invalid public key")
	ErrNotImplemented                   = statusError(CodeNotImplemented, "not implemented")
)

func statusError(code codes.Code, msg string) error {
	return status.New(code, msg).Err()
}
