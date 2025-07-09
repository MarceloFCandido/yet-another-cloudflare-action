package utils

import (
	"errors"
	"testing"
)

func TestPanicOnError(t *testing.T) {
	t.Run("should not panic when error is nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The code panicked when it should not have")
			}
		}()
		PanicOnError(nil)
	})

	t.Run("should panic when error is not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic when it should have")
			}
		}()
		PanicOnError(errors.New("test error"))
	})
}
