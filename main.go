package main

import (
	"log"
	"os"

	_ "example.com/matchers"
	"example.com/search"
)

func init() {
	log.SetOutput(os.Stdout)
}
func main() {
	if len(os.Args) <= 1 {
		log.Printf("[Usage]: %s <keyword>", os.Args[0])
		os.Exit(1)
	}
	search.Run(os.Args[1])
}
