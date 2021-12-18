package jsonlex

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

// bytes.Buffer stays bytes.Buffer
func TestEnsureAtLeastSingleByteUnreadableReader_1(t *testing.T) {
	given := new(bytes.Buffer)

	actual := EnsureAtLeastSingleByteUnreadableReader(given)

	if actual != given {
		t.Errorf("EnsureAtLeastSingleByteUnreadableReader() = %v, want: %v", actual, given)
	}
}

// uselessTestReader will be wrapped into *singleByteUnreadableReader
func TestEnsureAtLeastSingleByteUnreadableReader_2(t *testing.T) {
	given := &uselessTestReader{}

	actual := EnsureAtLeastSingleByteUnreadableReader(given)

	if sbur, ok := actual.(*singleByteUnreadableReader); !ok {
		t.Errorf("EnsureAtLeastSingleByteUnreadableReader() = %v, should be of: %v", actual, reflect.TypeOf(&singleByteUnreadableReader{}))
	} else if sbur.delegate != given {
		t.Errorf("EnsureAtLeastSingleByteUnreadableReader().delegate = %v, want: %v", sbur.delegate, given)
	}
}

func TestSingleByteUnreadableReader_Read(t *testing.T) {
	delegate := bytes.NewBuffer([]byte("0123456789"))
	instance := &singleByteUnreadableReader{delegate: delegate}

	steps := []struct {
		doUnread      bool
		amount        int
		expected      string
		expectedError error
		left          string
	}{{ //0
		amount:   2,
		expected: `01`,
		left:     `23456789`,
	}, { //1
		amount:   1,
		expected: `2`,
		left:     `3456789`,
	}, { //2
		doUnread: true,
		left:     `3456789`,
	}, { //3
		amount:   1,
		expected: `2`,
		left:     `3456789`,
	}, { //4
		amount:   3,
		expected: `345`,
		left:     `6789`,
	}, { //5
		doUnread: true,
		left:     `6789`,
	}, { //6
		doUnread:      true,
		expectedError: io.ErrShortBuffer,
		left:          `6789`,
	}, { //7
		amount:   3,
		expected: `567`,
		left:     `89`,
	}, { //8
		amount:   0,
		expected: ``,
		left:     `89`,
	}}

	for i, step := range steps {
		if step.doUnread {
			actualErr := instance.UnreadByte()
			if actualErr != step.expectedError {
				t.Errorf("%d: instance.UnreadByte() = %v, want: %v", i, actualErr, step.expectedError)
			}
		} else {
			actual := make([]byte, step.amount)
			actualN, actualErr := instance.Read(actual)
			if actualErr != step.expectedError {
				t.Errorf("%d: instance.Read() err = %v, want: %v", i, actualErr, step.expectedError)
			}
			if actualN != len(step.expected) {
				t.Errorf("%d: instance.Read() n = %d, want: %d", i, actualN, len(step.expected))
			}
			if string(actual) != step.expected {
				t.Errorf("%d: actual = %q, want: %q", i, actual, step.expected)
			}
		}
		if delegate.String() != step.left {
			t.Errorf("%d: left = %q, want: %q", i, delegate.String(), step.left)
		}
	}
}

type uselessTestReader struct{}

func (r *uselessTestReader) Read([]byte) (int, error) {
	panic("should never be called")
}
