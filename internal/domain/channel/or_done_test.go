package channel_test

import (
	"testing"

	"github.com/lgeorgieff/go-playground/internal/domain/channel"
	"github.com/stretchr/testify/assert"
)

func Test_OrDone_done(t *testing.T) {
	const size = 10

	ints := make(chan int)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for i := 0; i < size; i++ {
			ints <- i
		}
	}()

	orDone := channel.OrDone(ints, done)

	counter := 0
	for i := range orDone {
		assert.Equal(t, counter, i)
		counter++
	}
	assert.LessOrEqual(t, counter, size)
}

func Test_OrDone_close(t *testing.T) {
	const size = 10

	ints := make(chan int)
	done := make(chan struct{})

	go func() {
		defer close(ints)
		for i := 0; i < size; i++ {
			ints <- i
		}
	}()

	orDone := channel.OrDone(ints, done)

	counter := 0
	for i := range orDone {
		assert.Equal(t, counter, i)
		counter++
	}
	assert.Equal(t, size, counter)
}
