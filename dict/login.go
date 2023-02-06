package dict

import (
	"fmt"
	"strings"
)

type Login struct {
	Line

	sep string

	User   string
	Passwd string
}

func (login *Login) String() string {
	return fmt.Sprintf("%s:%s", login.User, login.Passwd)
}

func (login *Login) GetSep() string {
	return login.sep
}

func (login *Login) SetSep(sep string) {
	login.sep = sep
}

// MakeDefaultLogin make Login dict line
func MakeDefaultLogin(str string) (Line, error) {
	sep := ":"
	line := strings.Split(str, sep)
	if len(line) != 2 {
		return nil, fmt.Errorf("split line error")
	}
	return &Login{sep: sep, User: line[0], Passwd: line[1]}, nil
}
