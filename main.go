package main

import "fmt"

func main() {
	fmt.Println("Hello, SealScript!")
}

func Pure[A any](x A) []A {
	return []A{x}
}
