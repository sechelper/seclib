package async

import "time"

// Goroutine asynchronous processing of data in chan
func Goroutine(n int, c *chan any, input func(...any), process func(...any)) {
	done := make(chan struct{})

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case args := <-*c:
					process(args)
				case <-done:
					return
				default:
				}
			}
		}()
	}

	input(c)

	// clear chan
	for {
		if len(*c) == 0 {
			close(done)
			return
		}
		time.Sleep(500)
	}
}
