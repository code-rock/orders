package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func GetTestMessages(doSomething func([]byte)) {
	dir := "./test-data"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, file := range files {
			func(file string) {
				jsonFile, _ := os.Open(fmt.Sprintf("%s/%s", dir, file))
				defer jsonFile.Close()
				byteValue, _ := ioutil.ReadAll(jsonFile)
				doSomething(byteValue)
				time.Sleep(1 * time.Second)
			}(file.Name())
		}
	}
}
