package dict_test

import (
	"fmt"
	"github.com/sechelper/seclib/dict"
	"log"
)

func ExampleNewDict() {
	path := "user-pass.txt"

	counter, err := dict.Counter(path)
	if err != nil {
		log.Fatal(err)
	}

	// default dict
	dt := dict.NewDefaultDict()
	dt.MakeLine = dict.MakeDefaultLogin // default use dict.MakeDefaultStrLine

	// custom dict
	//dt := &dict.Dict{
	//	Lines:      make(chan dict.Line, 1000),
	//	Done:       make(chan struct{}, 1),
	//	MakeLine:   dict.MakeDefaultLogin,
	//}

	// custom dict
	//dt := dict.NewDict(1000, dict.MakeDefaultLogin)

	defer dt.Close()
	go func() {
		if err := dt.LoadText(path); err != nil {
			return
		}
		dt.Close()
	}()

	for i := 0; i < counter; i++ {
		select {
		case line, ok := <-dt.Lines:
			if !ok {
				return
			}
			login := line.(*dict.Login)
			fmt.Println(login.User, login.Passwd)
		case err := <-dt.Err:
			fmt.Println(err)
			return
		case <-dt.Done:
			return
		}
	}
}
