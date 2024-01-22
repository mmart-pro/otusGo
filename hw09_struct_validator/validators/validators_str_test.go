package validators

import (
	"fmt"
	"testing"

	"github.com/mmart-pro/otusGo/hw09structvalidator/errs"
	"github.com/stretchr/testify/require"
)

func TestValidatorStrCreateErr(t *testing.T) {
	cases := []struct {
		in  string
		err error
	}{
		{in: " len:200", err: errs.ErrInvalidValidator},
		{in: "len:x", err: errs.ErrLenStrValidatorInvalid},
		{in: "len:", err: errs.ErrLenStrValidatorInvalid},
		{in: " regexp:200", err: errs.ErrInvalidValidator},
		{in: "regexp:", err: errs.ErrRegexpValidatorInvalid},
		{in: "regexp:\\", err: errs.ErrRegexpCompilationError},
		{in: "max:200", err: errs.ErrInvalidValidator},
		{in: " in:200,20", err: errs.ErrInvalidValidator},
		{in: "in :foo,bar", err: errs.ErrInvalidValidator},
	}
	for _, v := range cases {
		v := v
		t.Run(fmt.Sprintf("case %s", v.in), func(t *testing.T) {
			x, err := NewStrValidator(v.in)
			require.Nil(t, x)
			require.ErrorIs(t, err, v.err)
		})
	}
}

func TestLenValidatorValid(t *testing.T) {
	x, err := NewStrValidator("len:3")
	require.NotNil(t, x)
	require.NoError(t, err)
	require.Nil(t, x.Valid("300"))
	require.Nil(t, x.Valid("abc"))
	require.ErrorIs(t, x.Valid("abcd"), errs.ErrLenValidation)
}

func TestRegexpValidatorValid(t *testing.T) {
	x, err := NewStrValidator("regexp:\\d+")
	require.NoError(t, err)
	require.NotNil(t, x)
	require.Nil(t, x.Valid("55"))
	require.Nil(t, x.Valid("a54b"))
	require.ErrorIs(t, x.Valid("abc"), errs.ErrRegexpValidation)
	x, err = NewStrValidator("regexp:\\d+\\b")
	require.NoError(t, err)
	require.Nil(t, x.Valid("54"))
	require.ErrorIs(t, x.Valid("a23bc"), errs.ErrRegexpValidation)
}

func TestInStrValidatorValid(t *testing.T) {
	x, err := NewStrValidator("in:beer")
	require.NoError(t, err)
	require.NotNil(t, x)
	require.Nil(t, x.Valid("beer"))
	require.ErrorIs(t, x.Valid("bear"), errs.ErrInStrValidation)
	require.ErrorIs(t, x.Valid(" beer"), errs.ErrInStrValidation)
	require.ErrorIs(t, x.Valid("beer "), errs.ErrInStrValidation)
	require.ErrorIs(t, x.Valid(" beer "), errs.ErrInStrValidation)
	x, err = NewStrValidator("in:foo,bar")
	require.NoError(t, err)
	require.NotNil(t, x)
	require.Nil(t, x.Valid("foo"))
	require.Nil(t, x.Valid("bar"))
	require.ErrorIs(t, x.Valid("201"), errs.ErrInStrValidation)
	require.ErrorIs(t, x.Valid(" foo"), errs.ErrInStrValidation)
	require.ErrorIs(t, x.Valid("bar "), errs.ErrInStrValidation)
}
