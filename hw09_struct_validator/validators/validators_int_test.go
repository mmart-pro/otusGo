package validators

import (
	"fmt"
	"testing"

	"github.com/mmart-pro/otusGo/hw09structvalidator/errs"
	"github.com/stretchr/testify/require"
)

func TestValidatorCreateErr(t *testing.T) {
	cases := []struct {
		in  string
		err error
	}{
		{in: "min:x", err: errs.ErrMinIntValidatorInvalid},
		{in: "min:", err: errs.ErrMinIntValidatorInvalid},
		{in: " min:200", err: errs.ErrInvalidValidator},
		{in: "max:x", err: errs.ErrMaxIntValidatorInvalid},
		{in: "max:", err: errs.ErrMaxIntValidatorInvalid},
		{in: " max:200", err: errs.ErrInvalidValidator},
		{in: "in:200,,20", err: errs.ErrInIntValidatorInvalid},
		{in: "in:200,s,20", err: errs.ErrInIntValidatorInvalid},
		{in: "in:", err: errs.ErrInIntValidatorInvalid},
		{in: " in:200,20", err: errs.ErrInvalidValidator},
	}
	for _, v := range cases {
		v := v
		t.Run(fmt.Sprintf("case %s", v.in), func(t *testing.T) {
			x, err := NewIntValidator(v.in)
			require.Nil(t, x)
			require.ErrorIs(t, err, v.err)
		})
	}
}

func TestMinValidatorValid(t *testing.T) {
	x, err := NewIntValidator("min:300")
	require.NotNil(t, x)
	require.Nil(t, err)
	require.Nil(t, x.Valid(300))
	require.Nil(t, x.Valid(301))
	require.ErrorIs(t, x.Valid(299), errs.ErrMinValidation)
}

func TestMaxValidatorValid(t *testing.T) {
	x, err := NewIntValidator("max:55")
	require.NotNil(t, x)
	require.Nil(t, err)
	require.Nil(t, x.Valid(55))
	require.Nil(t, x.Valid(54))
	require.ErrorIs(t, x.Valid(56), errs.ErrMaxValidation)
}

func TestInIntValidatorValid(t *testing.T) {
	x, err := NewIntValidator("in:200")
	require.NotNil(t, x)
	require.Nil(t, err)
	require.Nil(t, x.Valid(200))
	require.ErrorIs(t, x.Valid(201), errs.ErrInIntValidation)
	require.ErrorIs(t, x.Valid(199), errs.ErrInIntValidation)
	x, err = NewIntValidator("in:300,400,500")
	require.NotNil(t, x)
	require.Nil(t, err)
	require.Nil(t, x.Valid(300))
	require.Nil(t, x.Valid(400))
	require.Nil(t, x.Valid(500))
	require.ErrorIs(t, x.Valid(-300), errs.ErrInIntValidation)
	require.ErrorIs(t, x.Valid(0), errs.ErrInIntValidation)
	require.ErrorIs(t, x.Valid(600), errs.ErrInIntValidation)
}
