package test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

// assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...any) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\033[31m%s:%d: "+msg+"\033[39m\n\n",
			append([]any{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\033[31m%s:%d: unexpected error: %s\033[39m\n\n",
			filepath.Base(file),
			line,
			err.Error(),
		)
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act any) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n",
			filepath.Base(file),
			line,
			exp,
			act,
		)
		tb.FailNow()
	}
}

func SetPrivateField[T any](s T, fieldName string, value any) error {
	v := reflect.ValueOf(s)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("SetPrivateField: expected pointer to struct, got %T", s)
	}

	v = v.Elem()
	f := v.FieldByName(fieldName)

	if !f.IsValid() {
		return fmt.Errorf("SetPrivateField: no such field %q", fieldName)
	}

	reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
	return nil
}
