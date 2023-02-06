package dict

type StrLine string

func (sl *StrLine) GetSep() string {
	//TODO implement me
	panic("implement me")
}

func (sl *StrLine) SetSep(s string) {
	//TODO implement me
	panic("implement me")
}

func (sl *StrLine) String() string {
	return string(*sl)
}

func MakeDefaultStrLine(str string) (Line, error) {
	line := StrLine(str)
	return &line, nil
}
