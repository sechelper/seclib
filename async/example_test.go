package async_test

import (
	"fmt"
	"github.com/sechelper/seclib/async"
)

func ExampleGoroutine() {
	input := func(c *chan any) {
		for i := 0; i <= 100000; i++ {
			*c <- i
		}
	}

	process := func(p ...any) {
		fmt.Println(p[0])
	}

	async.Goroutine(30, input, process)

}
