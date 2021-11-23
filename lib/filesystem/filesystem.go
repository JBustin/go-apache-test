package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Filer interface {
	Exists(filepath string) bool
	Copy(src string, dst string) error
	Read(filepath string) ([]byte, error)
	Write(filepath string, content string) error
	List(dirname string, filterFn func(filepath string) bool) ([]string, error)
	Append(filepath string, content string) error
}

type filesystem struct{}

func New() filesystem {
	return filesystem{}
}

func (f filesystem) Exists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func (f filesystem) Copy(src string, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, input, 0644)
}

func (f filesystem) Read(filepath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filepath)
	return content, err
}

func (f filesystem) Write(filepath string, content string) error {
	return ioutil.WriteFile(filepath, []byte(content), 0644)
}

func (f filesystem) Append(filepath string, content string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.WriteString(content)

	return err
}

func (f filesystem) List(dirname string, filterFn func(filepath string) bool) ([]string, error) {
	var list []string
	err := filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsDir() && filterFn(path) {
				list = append(list, path)
			}
			return nil
		})

	return list, err
}
