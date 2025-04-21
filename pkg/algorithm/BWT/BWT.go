package BWT

import (
	"bytes"
	"fmt"
	"sort"
)

func Encode(data []byte) ([]byte, uint32) {
	// 1. Generate all rearrangements
	rotations := make([][]byte, len(data))
	for i := range rotations {
		rotations[i] = append(data[i:], data[:i]...)
	}

	// 2. Sort
	sort.Slice(rotations, func(i, j int) bool {
		return bytes.Compare(rotations[i], rotations[j]) < 0
	})

	// 3. Choose last column
	lastColumn := make([]byte, len(data))
	var originalIndex int
	for i, rot := range rotations {
		if bytes.Equal(rot, data) {
			originalIndex = i
		}
		lastColumn[i] = rot[len(rot)-1]
	}

	return lastColumn, uint32(originalIndex)
}

func Decode(lastColumn []byte, primaryIndex int) ([]byte, error) {
	if len(lastColumn) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	if primaryIndex < 0 || primaryIndex >= len(lastColumn) {
		return nil, fmt.Errorf("invalid primary index %d for length %d", primaryIndex, len(lastColumn))
	}

	n := len(lastColumn)

	// 1. Frequency
	count := make(map[byte]int)
	for _, b := range lastColumn {
		count[b]++
	}

	// 2.First column
	chars := make([]byte, 0, len(count))
	for c := range count {
		chars = append(chars, c)
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	// 3. Starting positions
	firstColStart := make(map[byte]int)
	total := 0
	for _, c := range chars {
		firstColStart[c] = total
		total += count[c]
	}

	// 4. LF-mapping
	lf := make([]int, n)
	symbolCount := make(map[byte]int)

	for i := 0; i < n; i++ {
		b := lastColumn[i]
		lf[i] = firstColStart[b] + symbolCount[b]
		symbolCount[b]++
	}

	// 5. Restore original string
	result := make([]byte, n)
	idx := primaryIndex

	for i := n - 1; i >= 0; i-- {
		if idx < 0 || idx >= len(lastColumn) {
			return nil, fmt.Errorf("invalid LF index %d at position %d", idx, i)
		}
		result[i] = lastColumn[idx]
		idx = lf[idx]
	}

	return result, nil
}
