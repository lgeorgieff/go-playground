package channel_test

import (
	"testing"

	"github.com/lgeorgieff/go-playground/internal/domain/channel"
	"github.com/stretchr/testify/assert"
)

func Test_Take_0(t *testing.T) {
	const size = 10

	done := make(chan struct{})
	stream := channel.Source(size, done, channel.DefaultIntGenerator)

	n := channel.TakeN(0, stream, done)
	v, ok := <-n
	assert.Zero(t, v)
	assert.False(t, ok)
}

func Test_Take_done(t *testing.T) {
	const size = 10

	done := make(chan struct{})
	stream := channel.Source(size, make(chan struct{}), channel.DefaultIntGenerator)

	n := channel.TakeN(size, stream, done)

	for i := 0; i < size; i++ {
		if i == size/2 {
			close(done)
			break
		}

		v, ok := <-n
		assert.Equal(t, i, v)
		assert.True(t, ok)
	}

	v, ok := <-n
	assert.Zero(t, v)
	assert.False(t, ok)
}

func Test_Take(t *testing.T) {
	testCases := []struct {
		Name          string
		SourceSize    int
		N             int
		ExpectedElems []int
	}{
		{Name: "take 1", SourceSize: 10, N: 1, ExpectedElems: []int{0}},
		{Name: "take 5", SourceSize: 10, N: 5, ExpectedElems: []int{0, 1, 2, 3, 4}},
		{Name: "take all", SourceSize: 10, N: 10, ExpectedElems: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{Name: "take more than available", SourceSize: 10, N: 11, ExpectedElems: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			done := make(chan struct{})
			src := channel.Source(testCase.SourceSize, done, channel.DefaultIntGenerator)
			n := channel.TakeN(testCase.N, src, done)

			for _, expectedElem := range testCase.ExpectedElems {
				actualElem, ok := <-n
				assert.True(t, ok)
				assert.Equal(t, expectedElem, actualElem)
			}
			actualElem, ok := <-n
			assert.False(t, ok)
			assert.Zero(t, actualElem)
		})
	}
}
