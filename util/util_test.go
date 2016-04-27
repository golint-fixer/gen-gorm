package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	t.Parallel()
	assert.Panics(t, func() {
		HandleErr(errors.New("some error"))
	}, "Calling handleErr() should panic")
}
