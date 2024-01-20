package validators

import (
	"strconv"
	"strings"

	"github.com/mmart-pro/otusGo/hw09structvalidator/errs"
)

// --------------------------------------------------------------------------------------------

type IntValidator interface {
	Valid(v int64) error
}

func NewIntValidator(str string) (IntValidator, error) {
	i := strings.Index(str, ":")
	if i <= 0 {
		return nil, errs.ErrInvalidValidator
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
		return nil, errs.ErrInvalidValidator
	}
}

// --------------------------------------------------------------------------------------------

type MaxIntValidator struct {
	max int64
}

func NewMaxIntValidator(str string) (*MaxIntValidator, error) {
	max, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, errs.ErrMaxIntValidatorInvalid
	}
	return &MaxIntValidator{max: max}, nil
}

func (maxv MaxIntValidator) Valid(v int64) error {
	if v > maxv.max {
		return errs.ErrMaxValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------

type MinIntValidator struct {
	min int64
}

func NewMinIntValidator(str string) (*MinIntValidator, error) {
	min, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, errs.ErrMinIntValidatorInvalid
	}
	return &MinIntValidator{min: min}, nil
}

func (minv MinIntValidator) Valid(v int64) error {
	if v < minv.min {
		return errs.ErrMinValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------

type InIntValidator struct {
	values []int64
}

func NewInIntValidator(str string) (*InIntValidator, error) {
	arr := strings.Split(str, ",")
	res := make([]int64, 0, len(arr))
	for _, v := range arr {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, errs.ErrInIntValidatorInvalid
		}
		res = append(res, val)
	}
	return &InIntValidator{values: res}, nil
}

func (validator InIntValidator) Valid(v int64) error {
	for _, a := range validator.values {
		if a == v {
			return nil
		}
	}
	return errs.ErrInIntValidation
}
