package channel

type Generator[T any] func(position int) T

var DefaultIntGenerator Generator[int] = func(position int) int {
	return position
}

func Source[T any](size int, done <-chan struct{}, fn Generator[T]) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		for i := 0; i < size; i++ {
			select {
			case <-done:
				return
			case stream <- fn(i):
			}
		}
	}()

	return stream
}
