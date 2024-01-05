package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorStrCreateErr(t *testing.T) {
	require.Nil(t, NewStrValidator("len:x"))
	require.Nil(t, NewStrValidator("len:"))
	require.Nil(t, NewStrValidator(" len:200"))
	//
	require.Nil(t, NewStrValidator("regexp:"))
	require.Nil(t, NewStrValidator(" regexp:200"))
	//
	require.Nil(t, NewStrValidator(" in:200,20"))
	require.Nil(t, NewStrValidator("in :foo,bar"))
}

func TestLenValidatorValid(t *testing.T) {
	x := NewStrValidator("len:3")
	require.NotNil(t, x)
	require.Nil(t, x.Valid("300"))
	require.Nil(t, x.Valid("abc"))
	require.ErrorIs(t, x.Valid("abcd"), ErrLenValidation)
}

func TestRegexpValidatorValid(t *testing.T) {
	x := NewStrValidator("regexp:\\d+")
	require.NotNil(t, x)
	require.Nil(t, x.Valid("55"))
	require.Nil(t, x.Valid("a54b"))
	require.ErrorIs(t, x.Valid("abc"), ErrRegexpValidation)
	x = NewStrValidator("regexp:\\d+\\b")
	require.Nil(t, x.Valid("54"))
	require.ErrorIs(t, x.Valid("a23bc"), ErrRegexpValidation)
}

func TestInStrValidatorValid(t *testing.T) {
	x := NewStrValidator("in:beer")
	require.NotNil(t, x)
	require.Nil(t, x.Valid("beer"))
	require.ErrorIs(t, x.Valid("bear"), ErrInStrValidation)
	require.ErrorIs(t, x.Valid(" beer"), ErrInStrValidation)
	require.ErrorIs(t, x.Valid("beer "), ErrInStrValidation)
	require.ErrorIs(t, x.Valid(" beer "), ErrInStrValidation)
	x = NewStrValidator("in:foo,bar")
	require.NotNil(t, x)
	require.Nil(t, x.Valid("foo"))
	require.Nil(t, x.Valid("bar"))
	require.ErrorIs(t, x.Valid("201"), ErrInStrValidation)
	require.ErrorIs(t, x.Valid(" foo"), ErrInStrValidation)
	require.ErrorIs(t, x.Valid("bar "), ErrInStrValidation)
}
