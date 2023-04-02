package main

import (
	"fmt"
	"liggi-go-jack-compiler/tokeniser"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No file specified")
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tokeniser := tokeniser.NewTokeniser(file)
	tokens, err := tokeniser.Tokenise()
	if err != nil {
		log.Fatal(err)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}
