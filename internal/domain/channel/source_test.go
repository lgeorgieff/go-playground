package channel_test

import (
	"testing"

	"github.com/lgeorgieff/go-playground/internal/domain/channel"
	"github.com/stretchr/testify/assert"
)

func Test_Source_0(t *testing.T) {
	done := make(chan struct{})
	src := channel.Source(0, done, channel.DefaultIntGenerator)

	v, ok := <-src
	assert.Zero(t, v)
	assert.False(t, ok)
}

func Test_Source(t *testing.T) {
	done := make(chan struct{})
	const size = 10
	gen := channel.Source(size, done, channel.DefaultIntGenerator)

	for i := 0; i < size; i++ {
		v, ok := <-gen
		assert.Equal(t, i, v)
		assert.True(t, ok)
	}
	v, ok := <-gen
	assert.Zero(t, v)
	assert.False(t, ok)
}

func Test_Generate_done(t *testing.T) {
	done := make(chan struct{})
	const size = 10
	gen := channel.Source(size, done, channel.DefaultIntGenerator)

	for i := 0; i < size; i++ {
		if i == 1 {
			close(done)
			break
		}

		v, ok := <-gen
		assert.Equal(t, i, v)
		assert.True(t, ok)
	}
	v, ok := <-gen
	assert.Zero(t, v)
	assert.False(t, ok)
}
