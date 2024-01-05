package validators

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	// - Для строк:
	ErrLenValidation    = errors.New("len validation failed")
	ErrRegexpValidation = errors.New("regexp validation failed")
	ErrInIntValidation  = errors.New("in int validation failed")
)

// --------------------------------------------------------------------------------------------

type StrValidator interface {
	Valid(v string) error
}

func NewStrValidator(str string) StrValidator {
	i := strings.Index(str, ":")
	if i <= 0 {
		return nil // unsupported validator
	}
	pref := str[:i]
	cond := str[i+1:]

	switch pref {
	case "len":
		return NewLenStrValidator(cond)
	case "regexp":
		return NewRegexpStrValidator(cond)
	case "in":
		return NewInStrValidator(cond)
	default:
		return nil // unsupported validator
	}
}

// --------------------------------------------------------------------------------------------

type LenStrValidator struct {
	len int
}

func NewLenStrValidator(str string) *LenStrValidator {
	reqLen, err := strconv.Atoi(str)
	if err != nil {
		return nil
	}
	return &LenStrValidator{len: reqLen}
}

func (lenv LenStrValidator) Valid(v string) error {
	if len(v) != lenv.len {
		return ErrLenValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------

type RegexpStrValidator struct {
	regexp *regexp.Regexp
}

func NewRegexpStrValidator(str string) *RegexpStrValidator {
	if str == "" {
		return nil
	}
	r, err := regexp.Compile(str)
	if err != nil {
		return nil
	}
	return &RegexpStrValidator{regexp: r}
}

func (validator RegexpStrValidator) Valid(v string) error {
	if !validator.regexp.MatchString(v) {
		return ErrRegexpValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------

type InStrValidator struct {
	values []string
}

func NewInStrValidator(str string) *InStrValidator {
	res := strings.Split(str, ",")
	return &InStrValidator{values: res}
}

func (validator InStrValidator) Valid(v string) error {
	for _, a := range validator.values {
		if a == v {
			return nil
		}
	}
	return ErrInStrValidation
}
