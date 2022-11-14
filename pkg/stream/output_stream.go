package stream

import "fmt"

var (
	OutputCh = make(chan string, 1000)
)

func Output() {
	for {
		line := <-OutputCh
		fmt.Println(line)
	}
}
