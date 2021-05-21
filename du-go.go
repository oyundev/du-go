package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var size int64

func dirSize(currPath string, info os.FileInfo) {
	defer wg.Done()
	files, _ := ioutil.ReadDir(currPath)

	for _, file := range files {
		if file.IsDir() {
			var newpath = fmt.Sprintf("%s/%s", currPath, file.Name())

			wg.Add(1)
			go dirSize(newpath, file)
		} else {
			atomic.AddInt64(&size, file.Size())
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: du-go <target_dir>\ndu-go c:\\temp")
	} else {
		var dir = os.Args[1]

		info, err := os.Lstat(dir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		wg.Add(1)
		go dirSize(dir, info)
		wg.Wait()

		fmt.Printf("%.2fM\t%s\n", float64(size)/(1024*1024), dir)
	}
}
