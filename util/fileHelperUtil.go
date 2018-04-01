package util

import (
	"log"
	"os"
)

//FileHelperUtil : helper methods for files
type FileHelperUtil struct {
	file []byte
}

//CheckDirFileExists : check file or dir exists
func (f *FileHelperUtil) CheckDirFileExists(path string) bool {

	_, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	}

	return true
}

//CreateDir : create directory
func (f *FileHelperUtil) CreateDir(dirpath string) {

	err := os.MkdirAll(dirpath, os.FileMode(0511))

	if err != nil {
		log.Println(err)
	}
}
