package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessTemplates(t *testing.T) {
	t.Parallel()
	var dir = "dist"

	// delete it if it exists
	if _, err := os.Stat(dir); err == nil {
		os.RemoveAll(dir)
	}
	processTemplates(database, "dist")

	_, err := os.Stat(dir)
	assert.NoError(t, err, "directory should exist")
	f, err := os.Stat(dir + "/models.go")
	assert.NoError(t, err, "models.go should exist")
	assert.True(t, f.Size() > 0, "file should not be empty")

	// clean up
	os.RemoveAll(dir)
}
