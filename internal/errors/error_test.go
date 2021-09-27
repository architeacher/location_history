package errors

import (
	"fmt"
	"reflect"
	"testing"
)

func TestKindString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id       string
		input    error
		expected string
	}{
		{
			"Test kind Other",
			E(Other),
			"unknown error kind",
		},
		{
			"Test kind Invalid",
			E(Invalid),
			"invalid input is provided",
		},
		{
			"Test kind NotFound",
			E(NotFound),
			"not found",
		},
		{
			"Test kind MinLength",
			E(MinLength),
			"not matching minimum length constraint",
		},
		{
			"Test kind Failure",
			E(Failure),
			"could not perform operation on the provided data %s",
		},
		{
			"Test kind Panic",
			E(Panic),
			"panic",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()
			assertEqual(t, testCase.expected, testCase.input.Error())
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id       string
		input    interface{}
		expected string
	}{
		{
			"Should append operation only",
			E(Operation("errors.TestError")),
			"errors.TestError",
		},
		{
			"Should append kind only",
			E(Invalid),
			"invalid entity is provided, entity can not be nil",
		},
		{
			"Should append error message only",
			E(fmt.Errorf("error message")),
			"error message",
		},
		{
			"Should append no error message",
			&Error{},
			"no error",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()
			assertEqual(t, testCase.expected, testCase.input.(*Error).Error())
		})
	}
}

func TestIs(t *testing.T) {
	t.Parallel()

	type (
		kindStruct struct {
			kind Kind
			err  error
		}
	)

	testCases := []struct {
		id       string
		input    kindStruct
		expected interface{}
	}{
		{
			"Nil errors should not have kind",
			kindStruct{
				Other,
				error(nil),
			},
			false,
		},
		{
			"Normal errors should not have kind",
			kindStruct{
				Invalid,
				fmt.Errorf("string %s", "error"),
			},
			false,
		},
		{
			"Should return true when the error have the same kind",
			kindStruct{
				MinLength,
				E(MinLength),
			},
			true,
		},
		{
			"Should return false when the error doesn't have the same kind",
			kindStruct{
				MinLength,
				E(NotFound),
			},
			false,
		},
		{
			"Should return true when previous errors should have the same kind",
			kindStruct{
				MinLength,
				E(Other, E(MinLength)),
			},
			true,
		},
		{
			"Should return false when the kind is Other",
			kindStruct{
				Failure,
				E(Other),
			},
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()
			assertEqual(t, testCase.expected, Is(testCase.input.kind, testCase.input.err))
		})
	}
}

func TestEWithNoArgs(t *testing.T) {
	t.Parallel()

	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("E() did not panic")
		}
	}()
	_ = E()
}

func TestE(t *testing.T) {
	testCases := []struct {
		id    string
		input Error
	}{
		{
			"Should set input arguments with Error object",
			Error{
				Op:   Operation("errors.TestE"),
				Kind: Panic,
				Err:  E(NotFound),
			},
		},
		{
			"Should set input arguments with normal error",
			Error{
				Op:   Operation("errors.TestE"),
				Kind: Panic,
				Err:  Errorf("error message"),
			},
		},
		{
			"Should set unknown error arguments",
			Error{
				Op:   Operation("errors.TestE"),
				Kind: Panic,
				Err:  nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			err := E(testCase.input.Op, testCase.input.Kind, testCase.input.Err)

			if testCase.input.Err == nil {
				assertEqual(t, "unknown type <nil>, value <nil> in error call", err.Error())
				return
			}

			errObject := err.(*Error)
			assertEqual(t, testCase.input.Op, errObject.Op)
			assertEqual(t, testCase.input.Kind, errObject.Kind)
			assertEqual(t, testCase.input.Err, errObject.Err)
		})
	}
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if (expected == nil || actual == nil) && expected != actual {
		reportError(t, expected, actual)
	}

	if !reflect.DeepEqual(expected, actual) {
		reportError(t, expected, actual)
	}
}

func reportError(t *testing.T, expected, actual interface{}) {
	t.Errorf("\nNot Equal:\nExpected: %v\nGot: %v", expected, actual)
}
