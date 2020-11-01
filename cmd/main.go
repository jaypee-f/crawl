package main

import (
	"fmt"
	"os"

	"github.com/jaypee-f/crawler/internal/consumer"
)

// TODO
// Add concurrency will need it later
// Split into packages and testing

func main() {
	fmt.Println("let's crawl")
	if len(os.Args) < 2 {
		fmt.Println("provide a url to crawl")
	}

	c := consumer.New(os.Args[1])

	c.Start()
}
