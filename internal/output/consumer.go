package output

import "fmt"

func StartJsonMessageParser(ch <-chan string) {
	for msg := range ch {
		fmt.Println("output:" + msg)
	}
}
