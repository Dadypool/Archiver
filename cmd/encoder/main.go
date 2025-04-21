package main

import (
	"fmt"
	"github.com/Dadypool/archiver/internal/encoder"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: decoder <infile> <zipfile>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	err := encoder.EncodeFile(inputFile, outputFile)
	if err != nil {
		fmt.Printf("Decoding error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File %s successfully decompressed to %s\n", inputFile, outputFile)
}
