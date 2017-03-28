package utils

import (
	"io/ioutil"
	"os"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func WriteFile(outputFilePath string, content string) error {
	return ioutil.WriteFile(outputFilePath, []byte(content), 0666)
}
