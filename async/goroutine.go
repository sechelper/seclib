package async

import (
	"sync"
)

// Goroutine asynchronous processing of data in chan
func Goroutine(n int, input func(*chan any), process func(...any)) {
	done := make(chan struct{})
	c := make(chan any, n)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case args, ok := <-c:
					if !ok {
						return
					}
					process(args)
				case <-done:
					return
				}
			}

		}()

	}

	input(&c)
	close(c)
	wg.Wait() // wait all goroutine done
}
