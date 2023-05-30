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
	Check(err)
	fmt.Println(len(files))
	for _, file := range files {
		fmt.Println(file)
		go func(file string) {
			data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file))
			Check(err)
			fmt.Print(data)

			jsonFile, err := os.Open(fmt.Sprintf("%s/%s", dir, file))
			Check(err)
			defer jsonFile.Close()
			byteValue, err := ioutil.ReadAll(jsonFile)
			Check(err)
			doSomething(byteValue)
			time.Sleep(time.Duration(rand.Int63n(1000)) * time.Second)
		}(file.Name())
	}
}
