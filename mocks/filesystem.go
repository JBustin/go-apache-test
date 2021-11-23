package mocks

type WriteArgs struct {
	Filepath string
	Content  string
}

type filesystem struct {
	Calls      map[string]int
	StubExists bool
	StubList   []string
	WriteArgs  []WriteArgs
	StubRead   func(filepath string) ([]byte, error)
	StubWrite  func(filepath string, content string) error
	StubAppend func(filepath string, content string) error
}

func NewFilesystem() filesystem {
	f := filesystem{}
	f.Reset()
	return f
}

func (f *filesystem) Reset() {
	f.Calls = map[string]int{
		"Exists": 0,
		"Copy":   0,
		"Read":   0,
		"Write":  0,
		"Append": 0,
		"List":   0,
	}
	f.Calls = map[string]int{}
	f.StubExists = true
	f.StubList = []string{}
	f.StubRead = func(filepath string) ([]byte, error) { return []byte{}, nil }
	f.StubWrite = func(filepath string, content string) error { return nil }
	f.StubAppend = func(filepath string, content string) error { return nil }
}

func (f *filesystem) SetStubRead(fn func(filepath string) ([]byte, error)) {
	f.StubRead = fn
}

func (f *filesystem) SetStubWrite(fn func(filepath string, content string) error) {
	f.StubWrite = fn
}

func (f *filesystem) SetStubAppend(fn func(filepath string, content string) error) {
	f.StubAppend = fn
}

func (f filesystem) Exists(filepath string) bool {
	f.Calls["Exists"]++
	return f.StubExists
}

func (f filesystem) Copy(src string, dst string) error {
	f.Calls["Copy"]++
	return nil
}

func (f filesystem) Read(filepath string) ([]byte, error) {
	f.Calls["Read"]++
	return f.StubRead(filepath)
}

func (f *filesystem) Write(filepath string, content string) error {
	f.Calls["Write"]++
	return f.StubWrite(filepath, content)
}

func (f filesystem) Append(filepath string, content string) error {
	f.Calls["Append"]++
	return f.StubAppend(filepath, content)
}

func (f filesystem) List(dirname string, filterFn func(filepath string) bool) ([]string, error) {
	f.Calls["List"]++
	return f.StubList, nil
}
