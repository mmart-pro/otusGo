package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mmart-pro/otusGo/hw09structvalidator/errs"
	"github.com/mmart-pro/otusGo/hw09structvalidator/validators"
)

type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s:%s", e.Field, e.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, e := range v {
		builder.WriteString(fmt.Sprintf("%s\n", e.Error()))
	}
	return builder.String()
}

func Validate(v interface{}) error {
	typ := reflect.TypeOf(v)
	kind := typ.Kind()
	if kind != reflect.Struct {
		return errs.ErrInvalidArgument
	}

	errs := ValidationErrors{}

	fieldsCount := typ.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		tagValidators := strings.Split(tag, "|")
		val := reflect.ValueOf(v).Field(i)
		// чекаем поле структуры на соответствие значения валидаторам
		fieldErrors, err := checkValue(tagValidators, val)
		if err != nil {
			return err
		}
		if fieldErrors == nil {
			continue
		}
		// вписать ошибки в общий массив ошибок
		for _, e := range fieldErrors {
			errs = append(errs, ValidationError{Field: field.Name, Err: e})
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func checkValue(tags []string, value reflect.Value) ([]error, error) {
	switch value.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return checkInt(tags, value.Int())
	case reflect.String:
		return checkStr(tags, value.String())
	case reflect.Slice, reflect.Array:
		if value.Len() == 0 {
			return nil, nil // => ok
		} else {
			// это ужасно просто
			switch value.Index(0).Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				result := make([]error, 0)
				for i := 0; i < value.Len(); i++ {
					el := value.Index(i)
					fieldErrors, err := checkInt(tags, el.Int())
					if err != nil {
						return nil, err
					}
					result = append(result, fieldErrors...)
				}
				return result, nil
			case reflect.String:
				result := make([]error, 0)
				for i := 0; i < value.Len(); i++ {
					el := value.Index(i)
					fieldErrors, err := checkStr(tags, el.String())
					if err != nil {
						return nil, err
					}
					result = append(result, fieldErrors...)
				}
				return result, nil
			default:
				return nil, errs.ErrUnsupportedType
			}
		}
	default:
		return nil, errs.ErrUnsupportedType
	}
}

func checkInt(tags []string, value int64) ([]error, error) {
	res := []error{}
	for _, t := range tags {
		v, err := validators.NewIntValidator(t)
		if err != nil {
			return nil, err
		}
		if err := v.Valid(value); err != nil {
			res = append(res, err)
		}
	}
	return res, nil
}

func checkStr(tags []string, value string) ([]error, error) {
	res := []error{}
	for _, t := range tags {
		v, err := validators.NewStrValidator(t)
		if err != nil {
			return nil, err
		}
		if err := v.Valid(value); err != nil {
			res = append(res, err)
		}
	}
	return res, nil
}
