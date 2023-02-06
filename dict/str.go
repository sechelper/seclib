package dict

type StrLine struct {
	Line

	Str string
}

func (sl StrLine) String() string {
	return sl.Str
}

func MakeDefaultStrLine(str string) (Line, error) {
	return StrLine{Str: str}, nil
}
