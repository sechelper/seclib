package dict_test

import (
	"fmt"
	"github.com/sechelper/seclib/dict"
	"log"
	"os"
)

func ExampleNewDict() {
	op, err := os.Open("users.txt")
	if err != nil {
		log.Fatal(err)
	}
	d := dict.NewDict(op)

	for d.Scan() {
		if line, err := d.Line(); err == nil {
			fmt.Println(line)
		}
	}
}

func ExampleNewDictForFile() {
	d, err := dict.NewDictForFile("users.txt")
	if err != nil {
		log.Fatal(err)
	}

	for d.Scan() {
		if line, err := d.Line(); err == nil {
			fmt.Println(line)
		}
	}
}

func ExampleLoginLineFunc() {
	d, err := dict.NewDictForFile("user-pass.txt")
	if err != nil {
		log.Fatal(err)
	}
	d.LineFunc(dict.LoginLineFunc)
	for d.Scan() {
		if line, err := d.Line(); err == nil {
			fmt.Println(line.(dict.LoginLine).User, line.(dict.LoginLine).Passwd)
		}

	}
}
