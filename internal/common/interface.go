package common

import "os"

type FileWriter[T any] interface {
	GetExtension() string
	Write(file *os.File, data T) error
}
