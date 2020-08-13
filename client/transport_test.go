package client

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	var target untrustedError
	untrustErr := untrustedError{}
	err := fmt.Errorf("wrapping err: %w", untrustErr)
	if !errors.As(err, &target) {
		t.Error("error did not unwrap to untrustedError")
	}
	if !Untrusted(err) {
		t.Error("err was not recognized as untrusted")
	}
}
