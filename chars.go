package main

import "fmt"

func main() {

	for ch := 0x02580; ch <= 0x258f; ch++ {
		fmt.Printf("hex:0x%x, dec:%d, char:%c\n", ch, ch, ch)
	}

}
