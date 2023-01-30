package async

// Goroutine asynchronous processing of data in chan
func Goroutine(n int, c *chan any, input func(...any), process func(...any)) {

	go input(c)

	for i := 0; i < n; i++ {
		go func() {
			for {
				process(<-*c)
			}
		}()
	}
}
