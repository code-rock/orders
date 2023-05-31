package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
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
			go func(file string) {
				jsonFile, _ := os.Open(fmt.Sprintf("%s/%s", dir, file))
				defer jsonFile.Close()
				byteValue, _ := ioutil.ReadAll(jsonFile)
				doSomething(byteValue)
				time.Sleep(time.Duration(rand.Int63n(1000)) * time.Second)
			}(file.Name())
		}
	}
}
