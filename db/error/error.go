package dberror

import "fmt"

type errorCode string

const (
	ENOTFOUND errorCode = "not_found"
	EALREADYEXISTS errorCode = "already_exists"
)

func New(code errorCode) error {
	return &fileDBError{code}
}

type fileDBError struct {
	code errorCode
}

func (e *fileDBError) Error() string {
	return fmt.Sprint(e.code)
}
