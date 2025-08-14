package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func readChunksFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("error", "error", err)
	}
	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			break
		}
		data = data[:n]
		fmt.Println("data : ", string(data))
	}
}

func readFileLinebyLine(path string) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("err", "err", err)
	}

	line := ""

	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			break
		}

		data = data[:n]

		i := bytes.IndexByte(data, 10) // 10 is representing the Line Feed ascii number

		if i != -1 {
			line += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read data: %s\n", line)
			line = ""
		}
		line += string(data)

	}
	if len(string(line)) != 0 {
		fmt.Printf("read data: %s\n", line)
	}
}

//
//func getLinesChannel(f io.ReadCloser) <-chan string {
//	out := make(chan string, 1)
//
//	go func() {
//		defer f.Close()
//		defer close()
//	}()
//
//	return out
//
//}
