package main

import (
	"fmt"
	"log"
)

func main() {
	result, err := Process("input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
