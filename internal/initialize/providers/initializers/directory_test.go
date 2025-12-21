package initializers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDirectoryInitializer(t *testing.T) {
	fs := afero.NewMemMapFs()
	initializer := NewDirectoryInitializer("some/dir")

	// Test Init
	err := initializer.Init(context.Background(), fs, fs, nil, nil)
	assert.NoError(t, err)

	exists, _ := afero.DirExists(fs, "some/dir")
	assert.True(t, exists)

	// Test IsSetup
	setup, err := initializer.IsSetup(fs, fs, nil)
	assert.NoError(t, err)
	assert.True(t, setup)

	// Test Path
	assert.Equal(t, "some/dir", initializer.Path())
}