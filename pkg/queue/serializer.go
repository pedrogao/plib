package queue

import (
	"fmt"
	"io/ioutil"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type (
	Serializer interface {
		Load(path string, val any) error

		Dump(path string, val any) error

		DumpFile(file *os.File, val any) error
	}

	JsonSerializer struct {
	}
)

func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

func (s *JsonSerializer) Load(path string, val any) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s err: %s", path, err)
	}

	err = jsoniter.Unmarshal(bytes, val)
	if err != nil {
		return fmt.Errorf("unmarshal err: %s", err)
	}
	return nil
}

func (s *JsonSerializer) Dump(path string, val any) error {
	bytes, err := jsoniter.Marshal(val)
	if err != nil {
		return fmt.Errorf("marshal err: %s", err)
	}

	err = ioutil.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write %s file err: %s", path, err)
	}
	return nil
}

func (s *JsonSerializer) DumpFile(file *os.File, val any) error {
	bytes, err := jsoniter.Marshal(val)
	if err != nil {
		return fmt.Errorf("marshal err: %s", err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("write %s file err: %s", file.Name(), err)
	}
	return nil
}
