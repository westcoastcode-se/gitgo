package db

type ContentDatabase interface {
	Read(path string, i interface{}) error
	Write(path string, i interface{}, message string) error
}
