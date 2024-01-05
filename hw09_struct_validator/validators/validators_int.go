package validators

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// - Для чисел:
	ErrMinValidation   = errors.New("min validation failed")
	ErrMaxValidation   = errors.New("max validation failed")
	ErrInStrValidation = errors.New("in str validation failed")
)

// --------------------------------------------------------------------------------------------------------------------------------

type IntValidator interface {
	Valid(v int64) error
}

func NewIntValidator(str string) IntValidator {
	i := strings.Index(str, ":")
	if i <= 0 {
		return nil // unsupported validator
	}
	pref := str[:i]
	cond := str[i+1:]

	switch pref {
	case "min":
		return NewMinIntValidator(cond)
	case "max":
		return NewMaxIntValidator(cond)
	case "in":
		return NewInIntValidator(cond)
	default:
		return nil // unsupported validator
	}
}

// --------------------------------------------------------------------------------------------------------------------------------

type MaxIntValidator struct {
	max int64
}

func NewMaxIntValidator(str string) *MaxIntValidator {
	max, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil
	}
	return &MaxIntValidator{max: max}
}

func (maxv MaxIntValidator) Valid(v int64) error {
	if v > maxv.max {
		return ErrMaxValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------------------

type MinIntValidator struct {
	min int64
}

func NewMinIntValidator(str string) *MinIntValidator {
	min, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil
	}
	return &MinIntValidator{min: min}
}

func (minv MinIntValidator) Valid(v int64) error {
	if v < minv.min {
		return ErrMinValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------------------

type InIntValidator struct {
	values []int64
}

func NewInIntValidator(str string) *InIntValidator {
	arr := strings.Split(str, ",")
	res := make([]int64, 0, len(arr))
	for _, v := range arr {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil
		}
		res = append(res, val)
	}
	return &InIntValidator{values: res}
}

func (validator InIntValidator) Valid(v int64) error {
	for _, a := range validator.values {
		if a == v {
			return nil
		}
	}
	return ErrInIntValidation
}
