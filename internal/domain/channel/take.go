package channel

func TakeN[T any](n int, stream <-chan T, done <-chan struct{}) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)

		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if !ok {
					return
				}
				select {
				case <-done:
				case result <- v:
				}
			}
		}
	}()

	return result
}
