package main

import (
	"bufio"
	"go.pkg.wesai.com/p/base_lib/log"
	"io"
	"os"
	"fmt"
)

func readlines(filename string) (lines []string) {

	f, err := os.Open(filename)
	if err != nil {
		log.DLogger().Fatal(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for true {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.DLogger().Fatal(err)
		}
		lines = append(lines, string(line))

	}

	return

}

func main() {
	filename := "/Users/zhangyongsheng/Downloads/wdjanss.txt"
	lines := readlines(filename)
	fmt.Println(lines[0])
}
