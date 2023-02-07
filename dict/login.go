package dict

import (
	"bytes"
	"fmt"
)

// LoginLine user login use this Line
type LoginLine struct {
	Line

	User   string
	Passwd string
}

// LoginLineFunc make LoginLine
func LoginLineFunc(b []byte) (Line, error) {
	// TODO fix bug, like use:r:pass, return user is "use", the user should be "use:r"
	if i := bytes.IndexByte(b, ':'); i >= 0 {
		// We have a full newline-terminated line.
		return LoginLine{User: string(b[:i]), Passwd: string(b[i+1:])}, nil
	}
	return nil, fmt.Errorf("%s error not found ':' ", b)
}
