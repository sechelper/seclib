package dict_test

import (
	"fmt"
	"github.com/sechelper/seclib/async"
	"github.com/sechelper/seclib/dict"
	"log"
)

func ExampleNewDict() {
	path := "user-pass.txt"

	counter, err := dict.Counter(path)
	if err != nil {
		log.Fatal(err)
	}
	dt := dict.NewDict(1000, dict.MakeDefaultLogin)
	//dt := dict.DefaultDict
	//dt.MakeLine = dict.MakeDefaultLogin // default use dict.MakeDefaultStrLine
	//dt := &dict.Dict{
	//	Lines:      make(chan dict.Line, 1000),
	//	Done:       make(chan struct{}, 1),
	//	Annotation: "#",
	//	MakeLine:   dict.MakeDefaultLogin,
	//}

	defer dt.Close()
	go func() {
		if err := dt.LoadText(path); err != nil {
			log.Fatal(err)
		}
	}()

	async.Goroutine(10, func(c *chan any) {
		for i := 0; i < counter; i++ {
			line := <-dt.Lines
			*c <- line
		}
	}, func(a ...any) {
		fmt.Println(a[0].(*dict.Login).User, a[0].(*dict.Login).Passwd)
	})
}
