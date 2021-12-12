package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type JsonContentDatabase struct {
	rootPath string
}

func (d *JsonContentDatabase) Read(path string, i interface{}) error {
	filename := filepath.Join(d.rootPath, path)
	if _, err := os.Stat(filename); err != nil {
		return err
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, i)
}

func (d *JsonContentDatabase) Write(path string, i interface{}, message string) error {
	filename := filepath.Join(d.rootPath, path)
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
