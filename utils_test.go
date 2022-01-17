package di_test

import (
	"testing"

	. "github.com/go-tk/di"
	"github.com/stretchr/testify/assert"
)

func TestFullFunctionName(t *testing.T) {
	assert.Equal(t, "github.com/go-tk/di_test.TestFullFunctionName", FullFunctionName(TestFullFunctionName))
}
