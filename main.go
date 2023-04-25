package main

import (
	codegenerator "liggi-go-jack-compiler/code-generator"
	"liggi-go-jack-compiler/parser"
	"liggi-go-jack-compiler/tokeniser"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No folder specified")
	}

	folderPath := os.Args[1]

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".jack") {
			processFile(path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func processFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tokeniser := tokeniser.NewTokeniser(file)
	tokens, err := tokeniser.Tokenise()
	if err != nil {
		log.Fatal(err)
	}

	parser := parser.NewParser(tokens)
	syntax, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	codeGenerator := codegenerator.NewCodeGenerator(syntax)
	generated, err := codeGenerator.Generate()
	if err != nil {
		log.Fatal(err)
	}

	vmFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".vm"
	vmFile, err := os.Create(vmFilePath)
	if err != nil {
		log.Fatalf("error creating file %s: %v", vmFilePath, err)
	}
	defer vmFile.Close()

	_, err = vmFile.WriteString(generated)
	if err != nil {
		log.Fatalf("error writing to file %s: %v", vmFilePath, err)
	}
}
