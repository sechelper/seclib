package async_test

import (
	"fmt"
	"github.com/sechelper/seclib/async"
	"sync"
)

func ExampleGoroutine() {
	wg := sync.WaitGroup{}
	wg.Add(100000)

	c := make(chan any, 10)
	input := func(...any) {
		for i := 0; i <= 100000; i++ {
			c <- i
		}

	}
	process := func(p ...any) {
		fmt.Println(p[0])
		wg.Done()
	}

	async.Goroutine(30, &c, input, process)

	wg.Wait()
}
