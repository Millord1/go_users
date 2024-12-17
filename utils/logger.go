package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type Logger struct {
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
}

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func GetLogger(fileName string) *Logger {
	file := getOrCreateFile(fileName)
	return &Logger{
		Warning: log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:    log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error:   log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func getOrCreateFile(fileName string) *os.File {
	// Hack for /logs path
	// TODO: clean it
	path := basepath + "/../logs/" + fileName
	fmt.Println(path)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		f, createErr := os.Create(path)
		check(createErr)
		f.Close()
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	/* defer file.Close() */

	return file
}
