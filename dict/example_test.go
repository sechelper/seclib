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
		line, _ := d.Line()
		fmt.Println(line)
	}
}

func ExampleNewDictForFile() {
	d, err := dict.NewDictForFile("users.txt")
	if err != nil {
		log.Fatal(err)
	}

	for d.Scan() {
		line, _ := d.Line()
		fmt.Println(line)
	}
}

func ExampleLoginLineFunc() {
	d, err := dict.NewDictForFile("user-pass.txt")
	if err != nil {
		log.Fatal(err)
	}
	d.LineFunc(dict.LoginLineFunc)
	for d.Scan() {
		line, err := d.Line()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(line.(dict.LoginLine).User, line.(dict.LoginLine).Passwd)
	}
}
