package testify

import (
	"testing"

	"gitee.com/quant1x/pkg/testify/assert"
)

func TestImports(t *testing.T) {
	if assert.Equal(t, 1, 1) != true {
		t.Error("Something is wrong.")
	}
}
