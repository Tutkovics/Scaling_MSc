package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting...")

	algorithm(3000)

	fmt.Println("Bye my Friend!")
}

func tightAlgorithm(milliseconds int) {
	now := time.Now() //.Nanosecond() / 1000000

	duration := milliseconds * int(time.Millisecond)
	iteration := 0

	fmt.Println(time.Now())

	for now.Add(time.Duration(duration)).After(time.Now()) {
		iteration += 1
	}

	fmt.Println("# of runned iteration >", iteration, "<")

	fmt.Println(time.Now())

}
