package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mmart-pro/otusGo/hw09structvalidator/validators"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestErrInvalidArgument(t *testing.T) {
	err := Validate(0)
	require.ErrorIs(t, err, ErrInvalidArgument)
	err = Validate("0")
	require.ErrorIs(t, err, ErrInvalidArgument)
	err = Validate([]int{})
	require.ErrorIs(t, err, ErrInvalidArgument)
	err = Validate([]struct{}{{}, {}})
	require.ErrorIs(t, err, ErrInvalidArgument)
	err = Validate(struct{}{})
	require.NoError(t, err)
}

func TestErrInvalidValidator(t *testing.T) {
	err := Validate(struct {
		Field int `validate:"error:500"`
	}{
		Field: 200,
	})
	require.ErrorContains(t, err, ErrInvalidValidator.Error())
}

func TestErrUnsupportedType(t *testing.T) {
	err := Validate(struct {
		Field float32 `validate:"min:200"`
	}{
		Field: 0,
	})
	require.ErrorContains(t, err, ErrUnsupportedType.Error())
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// Response
		{
			in:          Response{Code: 200, Body: "somebody"},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 202},
			expectedErr: ValidationErrors{ValidationError{Field: "Code", Err: validators.ErrInIntValidation}},
		},
		{
			in:          Response{},
			expectedErr: ValidationErrors{ValidationError{Field: "Code", Err: validators.ErrInIntValidation}},
		},
		// Token
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in:          Token{Header: []byte{}, Payload: []byte{0, 1, 2}, Signature: make([]byte, 0)},
			expectedErr: nil,
		},
		// App
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in:          App{},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: validators.ErrLenValidation}},
		},
		{
			in:          App{Version: "1234"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: validators.ErrLenValidation}},
		},
		{
			in:          App{Version: "123456"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: validators.ErrLenValidation}},
		},
		// User
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    45,
				Email:  "none@some.com",
				Role:   "stuff",
				Phones: []string{"79972879197", "01234567890"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456123",
				Age:    45,
				Email:  "none@some.com",
				Role:   "stuff",
				Phones: []string{"79972879197"},
			},
			expectedErr: ValidationErrors{
				ValidationError{"ID", validators.ErrLenValidation},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456123",
				Age:    5,
				Email:  "none@some.com",
				Role:   "admin",
				Phones: []string{"79972879197"},
			},
			expectedErr: ValidationErrors{
				ValidationError{"ID", validators.ErrLenValidation},
				ValidationError{"Age", validators.ErrMinValidation},
			},
		},
		{
			in: User{
				ID:     "201",
				Age:    101,
				Email:  "none",
				Role:   "user",
				Phones: []string{"123", "321"},
			},
			expectedErr: ValidationErrors{
				ValidationError{"ID", validators.ErrLenValidation},
				ValidationError{"Age", validators.ErrMaxValidation},
				ValidationError{"Email", validators.ErrRegexpValidation},
				ValidationError{"Role", validators.ErrInStrValidation},
				ValidationError{"Phones", validators.ErrLenValidation},
				ValidationError{"Phones", validators.ErrLenValidation},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %T (%d)", tt.in, i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)

				var validationErrors ValidationErrors
				require.ErrorAs(t, err, &validationErrors)

				require.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
