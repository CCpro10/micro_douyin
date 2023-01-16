package util

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestFilter(t *testing.T) {
	InitFilter()

	ok, w := Filter.FindIn("你是不是傻逼")
	assert.Equal(t, true, ok)
	assert.Equal(t, "傻逼", w)
}

func BenchmarkFilter(b *testing.B) {
	InitFilter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Filter.FindIn("你是不是傻逼")
	}
}
