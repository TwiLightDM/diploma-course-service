package errs

import "errors"

var (
	ErrRecordNotFound              = errors.New("record not found")
	ErrDuplicateKey                = errors.New("duplicate key")
	ErrModuleIdAndCourseIdNotFound = errors.New("module id and course id not found")
)
