package utils

import (
	"os"
	"path"
)

func MakeDirAll(fileName string) (err error) {
	p := path.Dir(path.Clean(fileName))
	if _, err = os.Stat(p); os.IsNotExist(err) {
		err = os.MkdirAll(p, 0755)
	}
	return err
}
