package Huffman

import (
	"fmt"
	"sort"
)

type Code struct {
	Symbol byte
	Length uint8
}

func MapToCodes(codes map[byte]string) []Code {
	canonical := make([]Code, 0, len(codes))
	for symbol, code := range codes {
		canonical = append(canonical, Code{
			Symbol: symbol,
			Length: uint8(len(code)),
		})
	}

	sort.Slice(canonical, func(i, j int) bool {
		if canonical[i].Length == canonical[j].Length {
			return canonical[i].Symbol < canonical[j].Symbol
		}
		return canonical[i].Length < canonical[j].Length
	})

	return canonical
}

func CodesToMap(codes []Code) (map[byte]string, error) {
	result := make(map[byte]string)
	var currentCode uint32 = 0
	var prevLength uint8 = 0

	for _, item := range codes {
		if item.Length == 0 {
			return nil, fmt.Errorf("invalid code length 0 for symbol %d", item.Symbol)
		}

		if item.Length > prevLength {
			currentCode <<= item.Length - prevLength
		}

		result[item.Symbol] = fmt.Sprintf("%0*b", item.Length, currentCode)

		currentCode++
		prevLength = item.Length
	}

	return result, nil
}
