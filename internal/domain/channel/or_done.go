package channel

func OrDone[T any](stream <-chan T, done <-chan struct{}) <-chan T {
	orStream := make(chan T)

	go func() {
		defer close(orStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if !ok {
					return
				}
				select {
				case <-done:
				case orStream <- v:
				}
			}
		}
	}()

	return orStream
}
