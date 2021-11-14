package di_test

import (
	"testing"

	. "github.com/go-tk/di"
	"github.com/stretchr/testify/assert"
)

func TestPackagePath(t *testing.T) {
	assert.Equal(t, "github.com/go-tk/di_test", PackagePath(0))
}

func TestPackageName(t *testing.T) {
	assert.Equal(t, "di_test", PackageName(0))
}
