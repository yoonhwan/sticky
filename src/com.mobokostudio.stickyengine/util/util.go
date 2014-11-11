package util

import (
	"fmt"
	"os"
	
)

func GetFileData(filename string) (str string) {
	file, err := os.Open(filename)
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()
	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return
	}
	str = string(bs)
	return
}

func WriteFile(filename, data string) {
	file, err := os.Create(filename)
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()
	file.WriteString(data)
}

func PrintDir(filepath string)() {
	dir, err := os.Open(filepath)
    if err != nil {
        return
    }
    defer dir.Close()

    fileInfos, err := dir.Readdir(-1)
    if err != nil {
        return
    }
    for _, fi := range fileInfos {
        fmt.Println(fi.Name())
    }
}
