package validators

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/mmart-pro/otusGo/hw09structvalidator/errs"
)

// --------------------------------------------------------------------------------------------

type StrValidator interface {
	Valid(v string) error
}

func NewStrValidator(str string) (StrValidator, error) {
	i := strings.Index(str, ":")
	if i <= 0 {
		return nil, errs.ErrInvalidValidator
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
		return nil, errs.ErrInvalidValidator
	}
}

// --------------------------------------------------------------------------------------------

type LenStrValidator struct {
	len int
}

func NewLenStrValidator(str string) (*LenStrValidator, error) {
	reqLen, err := strconv.Atoi(str)
	if err != nil {
		return nil, errs.ErrLenStrValidatorInvalid
	}
	return &LenStrValidator{len: reqLen}, nil
}

func (lenv LenStrValidator) Valid(v string) error {
	if len(v) != lenv.len {
		return errs.ErrLenValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------

type RegexpStrValidator struct {
	regexp *regexp.Regexp
}

func NewRegexpStrValidator(str string) (*RegexpStrValidator, error) {
	if str == "" {
		return nil, errs.ErrRegexpValidatorInvalid
	}
	r, err := regexp.Compile(str)
	if err != nil {
		return nil, errs.ErrRegexpCompilationError
	}
	return &RegexpStrValidator{regexp: r}, nil
}

func (validator RegexpStrValidator) Valid(v string) error {
	if !validator.regexp.MatchString(v) {
		return errs.ErrRegexpValidation
	}
	return nil
}

// --------------------------------------------------------------------------------------------

type InStrValidator struct {
	values []string
}

func NewInStrValidator(str string) (*InStrValidator, error) {
	if str == "" {
		return nil, errs.ErrInStrValidatorInvalid
	}
	res := strings.Split(str, ",")
	return &InStrValidator{values: res}, nil
}

func (validator InStrValidator) Valid(v string) error {
	for _, a := range validator.values {
		if a == v {
			return nil
		}
	}
	return errs.ErrInStrValidation
}
