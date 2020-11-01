package main

import (
	"fmt"
	"os"

	"github.com/jaypee-f/crawler/internal/consumer"
)

func main() {
	fmt.Println("let's crawl")
	if len(os.Args) < 2 {
		fmt.Println("provide a url to crawl")
		os.Exit(-1)
	}

	c := consumer.New(os.Args[1])

	c.Start()
}
