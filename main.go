package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args)-1 < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args)-1 > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	if len(args)-1 == 1 {
		fmt.Printf("starting crawl of: %s\n", args[1])
	}

	page, err := getHTML(args[1])
	if err != nil {
		log.Fatalf("Error getting HTML: %v", err)
	}
	fmt.Printf("%v", page)
}
