package errs

import (
	"errors"
)

var (
	// - Для строк:
	ErrLenValidation    = errors.New("len validation failed")
	ErrRegexpValidation = errors.New("regexp validation failed")
	ErrInIntValidation  = errors.New("in int validation failed")
	// - Для чисел:
	ErrMinValidation   = errors.New("min validation failed")
	ErrMaxValidation   = errors.New("max validation failed")
	ErrInStrValidation = errors.New("in str validation failed")
)

// Ошибки программиста
var (
	ErrInvalidArgument  = errors.New("invalid argument, struct expected")
	ErrUnsupportedType  = errors.New("invalid or unsupported type")
	ErrInvalidValidator = errors.New("invalid or unsupported validator")

	ErrLenStrValidatorInvalid = errors.New("invalid value for string len validator")
	ErrRegexpValidatorInvalid = errors.New("invalid value for reexp validator")
	ErrRegexpCompilationError = errors.New("compiling rexep expression error")
	ErrInStrValidatorInvalid  = errors.New("invalid value for string in validator")

	ErrMaxIntValidatorInvalid = errors.New("invalid value for int max validator")
	ErrMinIntValidatorInvalid = errors.New("invalid value for int min validator")
	ErrInIntValidatorInvalid  = errors.New("invalid value for int in validator")
)
