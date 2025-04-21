package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the filename to compare")
	}

	originalFile := os.Args[1]
	encodedFile := originalFile + "_encoded"
	decodedFile := originalFile + "_decoded"

	originalData, err := ioutil.ReadFile(originalFile)
	if err != nil {
		log.Fatalf("Error reading original file: %v", err)
	}

	decodedData, err := ioutil.ReadFile(decodedFile)
	if err != nil {
		log.Fatalf("Error reading decoded file: %v", err)
	}

	if len(originalData) != len(decodedData) {
		fmt.Printf("Files %s have different sizes\n", originalFile)
		return
	}

	for i := 0; i < len(originalData); i++ {
		if originalData[i] != decodedData[i] {
			fmt.Printf("Files differ at byte %d\n", i)
			return
		}
	}

	fmt.Printf("Files %s and decoded one are identical\n", originalFile)

	originalFileInfo, err := os.Stat(originalFile)
	if err != nil {
		log.Fatalf("Error getting original file info: %v", err)
	}

	encodedFileInfo, err := os.Stat(encodedFile)
	if err != nil {
		log.Fatalf("Error getting encoded file info: %v", err)
	}

	fmt.Printf("Size of '%s': %d bytes\n", originalFile, originalFileInfo.Size())
	fmt.Printf("Size of '%s': %d bytes\n", encodedFile, encodedFileInfo.Size())

	fmt.Printf("Average bits for symbol: %.4f\n", float64(encodedFileInfo.Size()*8)/float64(originalFileInfo.Size()))

	fmt.Printf("H(X) = %.4f bits/symbol\n", entropy(originalData))
	fmt.Printf("H(X|X) = %.4f bits\n", conditionalEntropy1(originalData))
	fmt.Printf("H(X|XX) = %.4f bits\n", conditionalEntropy2(originalData))
}

func entropy(data []byte) float64 {
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}

	total := float64(len(data))
	var h float64
	for _, count := range freq {
		p := float64(count) / total
		h += -p * math.Log2(p)
	}
	return h
}

func conditionalEntropy1(data []byte) float64 {
	if len(data) < 2 {
		return 0
	}

	singleFreq := make(map[byte]int)
	pairFreq := make(map[[2]byte]int)
	total := len(data) - 1

	for i := 0; i < total; i++ {
		a, b := data[i], data[i+1]
		singleFreq[a]++
		pairFreq[[2]byte{a, b}]++
	}

	var h float64
	for pair, count := range pairFreq {
		a := pair[0]
		pAB := float64(count) / float64(total)
		pA := float64(singleFreq[a]) / float64(total)
		h += -pAB * math.Log2(pAB/pA)
	}

	return h
}

func conditionalEntropy2(data []byte) float64 {
	if len(data) < 3 {
		return 0
	}

	prefixFreq := make(map[[2]byte]int)
	tripleFreq := make(map[[3]byte]int)
	total := len(data) - 2

	for i := 0; i < total; i++ {
		a, b, c := data[i], data[i+1], data[i+2]
		prefix := [2]byte{a, b}
		triple := [3]byte{a, b, c}
		prefixFreq[prefix]++
		tripleFreq[triple]++
	}

	var h float64
	for triple, count := range tripleFreq {
		prefix := [2]byte{triple[0], triple[1]}
		pABC := float64(count) / float64(total)
		pAB := float64(prefixFreq[prefix]) / float64(total)
		h += -pABC * math.Log2(pABC/pAB)
	}

	return h
}
