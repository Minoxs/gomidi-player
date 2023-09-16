// Package utils contains useful functions to hide away
// some ugly code in a function call.
package utils

import (
	"fmt"
)

// IgnoreError will call function f and ignore error return.
// Useful to explicitly ignore errors in deferred functions.
func IgnoreError(f func() error) {
	_ = f()
}

// WrapOnError will return nil if errInner is nil.
// Otherwise returns an error that wraps both errInner and errOuter.
func WrapOnError(errInner error, errOuter error) error {
	if errInner != nil {
		return fmt.Errorf("%w: %w", errOuter, errInner)
	}
	return nil
}
