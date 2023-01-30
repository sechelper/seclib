package async

import "time"

// Goroutine asynchronous processing of data in chan
func Goroutine(n int, input func(*chan any), process func(...any)) {
	done := make(chan struct{})
	c := make(chan any, n)

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case args := <-c:
					process(args)
				case <-done:
					return
				}
			}
		}()
	}

	input(&c)

	// clear chan
	for {
		if len(c) == 0 {
			close(done)
			return
		}
		time.Sleep(500)
	}
}
