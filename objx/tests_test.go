package objx_test

import (
	"testing"

	"gitee.com/quant1x/pkg/objx"
	"gitee.com/quant1x/pkg/testify/assert"
)

func TestHas(t *testing.T) {
	m := objx.Map(TestMap)

	assert.True(t, m.Has("name"))
	assert.True(t, m.Has("address.state"))
	assert.True(t, m.Has("numbers[4]"))

	assert.False(t, m.Has("address.state.nope"))
	assert.False(t, m.Has("address.nope"))
	assert.False(t, m.Has("nope"))
	assert.False(t, m.Has("numbers[5]"))

	m = nil

	assert.False(t, m.Has("nothing"))
}
