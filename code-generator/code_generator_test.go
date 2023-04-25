package codegenerator

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"liggi-go-jack-compiler/parser"
	"liggi-go-jack-compiler/tokeniser"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func loadTestCases(path string) ([]string, error) {
	var testCases []string
	err := filepath.Walk(path, func(currPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && currPath != path {
			testCases = append(testCases, currPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return testCases, nil
}

func findJackFiles(path string) ([]string, error) {
	var jackFiles []string
	err := filepath.Walk(path, func(currPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jack") {
			jackFiles = append(jackFiles, currPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return jackFiles, nil
}

func generateVMFileFromJack(jackPath string) (string, error) {
	file, err := os.Open(jackPath)
	if err != nil {
		return "", fmt.Errorf("failed to open .jack file %s: %v", jackPath, err)
	}

	tokeniser := tokeniser.NewTokeniser(file)
	tokens, tokeniserErr := tokeniser.Tokenise()
	if tokeniserErr != nil {
		return "", fmt.Errorf("tokeniser error: %v", tokeniserErr)
	}

	parser := parser.NewParser(tokens)
	syntax, parserErr := parser.Parse()
	if parserErr != nil {
		return "", fmt.Errorf("parser error: %v", parserErr)
	}

	codeGenerator := NewCodeGenerator(syntax)
	generated, generatorErr := codeGenerator.Generate()
	if generatorErr != nil {
		return "", fmt.Errorf("code generator error: %v", generatorErr)
	}

	return generated, nil
}

func TestJackFiles(t *testing.T) {
	testCaseDirs, err := loadTestCases("../test-cases")
	if err != nil {
		t.Fatalf("failed to load test cases: %v", err)
	}

	for _, testCaseDir := range testCaseDirs {
		jackFiles, err := findJackFiles(testCaseDir)
		if err != nil {
			t.Fatalf("failed to find .jack files in %s: %v", testCaseDir, err)
		}

		for _, jackPath := range jackFiles {
			vmPath := strings.TrimSuffix(jackPath, ".jack") + ".vm"

			if _, err := os.Stat(vmPath); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				} else {
					t.Fatalf("error checking .vm file %s: %v", vmPath, err)
				}
			}

			fmt.Println("testing .jack file", jackPath)

			if _, err := os.Stat(vmPath); err == nil {
				expectedVMBytes, err := ioutil.ReadFile(vmPath)
				if err != nil {
					t.Fatalf("failed to read .vm file %s: %v", vmPath, err)
				}

				generatedVM, err := generateVMFileFromJack(jackPath)
				if err != nil {
					t.Fatalf("failed to process .jack file %s: %v", jackPath, err)
				}

				expectedVM := string(expectedVMBytes)

				if expectedVM != generatedVM {
					expectedVMSlice := strings.Split(expectedVM, "\n")
					generatedVMSlice := strings.Split(generatedVM, "\n")

					if diff := cmp.Diff(expectedVMSlice, generatedVMSlice); diff != "" {
						t.Errorf("mismatch in .vm file %s (-expected +got):\n%s", vmPath, diff)
					}
				} else {
					fmt.Println(".jack file", jackPath, "generated .vm file", vmPath, "successfully")
				}
			}
		}
	}
}
