package rad

import (
	"bytes"
)

// Comparable defines an interface for art values
type Comparable interface {
	// EqualTo returns true if the current value is equal to the provided value
	EqualTo(v interface{}) bool
}

// Bytes defines a byteslice that satisfies the comparable interface
type Bytes []byte

// EqualTo returns true if the compared value is equal to the target value
func (b Bytes) EqualTo(v interface{}) bool {
	cv, ok := v.(Bytes)
	if !ok {
		return false
	}

	return bytes.Equal(b, cv)
}

// String defines a string that satisfies the comparable interface
type String string

// EqualTo returns true if the compared value is equal to the target value
func (s String) EqualTo(v interface{}) bool {
	cv, ok := v.(String)
	if !ok {
		return false
	}

	return s == cv
}
