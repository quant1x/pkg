package objx_test

import (
	"testing"

	"gitee.com/quant1x/pkg/objx"
	"gitee.com/quant1x/pkg/testify/assert"
)

func TestHashWithKey(t *testing.T) {
	assert.Equal(t, "0ce84d8d01f2c7b6e0882b784429c54d280ea2d9", objx.HashWithKey("abc", "def"))
}
