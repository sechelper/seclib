package dict

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// Line dict line data
type Line interface {
	String() string
	GetSep() string
	SetSep(string)
}

// Dict data dictionary
type Dict struct {
	Lines chan Line // dict line channel
	Done  chan struct{}

	MakeLine func(string) (Line, error) // make dict Line
}

// NewDefaultDict default dict
// line chan size 1000 , annotation use '#', MakeLine use MakeDefaultStrLine
func NewDefaultDict() *Dict {
	return &Dict{
		Lines:    make(chan Line, 10),
		Done:     make(chan struct{}, 1),
		MakeLine: MakeDefaultStrLine,
	}
}

func NewDict(size int, makeLine func(string) (Line, error)) *Dict {
	return &Dict{
		Lines:    make(chan Line, size),
		Done:     make(chan struct{}, 1),
		MakeLine: makeLine,
	}
}

// LoadText load text dict from file
func (dict *Dict) LoadText(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		return err
	}

	for scanner.Scan() {
		if line, err := dict.MakeLine(scanner.Text()); err == nil {
			dict.Lines <- line
		} else {
			return err
		}
	}

	return nil
}

func (dict *Dict) Close() {
	dict.Done <- struct{}{}
	close(dict.Lines)
}

// Counter file line count
func Counter(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	r := bufio.NewReader(file)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
