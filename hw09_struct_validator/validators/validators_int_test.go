package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorCreateErr(t *testing.T) {
	require.Nil(t, NewIntValidator("min:x"))
	require.Nil(t, NewIntValidator("min:"))
	require.Nil(t, NewIntValidator(" min:200"))
	//
	require.Nil(t, NewIntValidator("max:x"))
	require.Nil(t, NewIntValidator("max:"))
	require.Nil(t, NewIntValidator(" max:200"))
	//
	require.Nil(t, NewIntValidator("in:200,,20"))
	require.Nil(t, NewIntValidator("in:200,s,20"))
	require.Nil(t, NewIntValidator("in:"))
	require.Nil(t, NewIntValidator(" in:200,20"))
}

func TestMinValidatorValid(t *testing.T) {
	x := NewIntValidator("min:300")
	require.NotNil(t, x)
	require.Nil(t, x.Valid(300))
	require.Nil(t, x.Valid(301))
	require.ErrorIs(t, x.Valid(299), ErrMinValidation)
}

func TestMaxValidatorValid(t *testing.T) {
	x := NewIntValidator("max:55")
	require.NotNil(t, x)
	require.Nil(t, x.Valid(55))
	require.Nil(t, x.Valid(54))
	require.ErrorIs(t, x.Valid(56), ErrMaxValidation)
}

func TestInIntValidatorValid(t *testing.T) {
	x := NewIntValidator("in:200")
	require.NotNil(t, x)
	require.Nil(t, x.Valid(200))
	require.ErrorIs(t, x.Valid(201), ErrInIntValidation)
	require.ErrorIs(t, x.Valid(199), ErrInIntValidation)
	x = NewIntValidator("in:300,400,500")
	require.NotNil(t, x)
	require.Nil(t, x.Valid(300))
	require.Nil(t, x.Valid(400))
	require.Nil(t, x.Valid(500))
	require.ErrorIs(t, x.Valid(-300), ErrInIntValidation)
	require.ErrorIs(t, x.Valid(0), ErrInIntValidation)
	require.ErrorIs(t, x.Valid(600), ErrInIntValidation)
}
