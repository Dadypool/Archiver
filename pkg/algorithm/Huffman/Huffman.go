package Huffman

import (
	"fmt"
	"sort"
)

type Node struct {
	Left, Right *Node
	Symbol      byte
	Frequency   int
}

func Encode(data []byte) ([]byte, []Code, error) {
	// 1. Build Tree
	tree := buildHuffmanTree(data)

	// 2. Generate raw codes
	rawCodes := make(map[byte]string)
	generateCodes(tree, "", rawCodes)

	// 3. Transform to canonical form
	canonicalSet := MapToCodes(rawCodes)

	canonicalMap, err := CodesToMap(canonicalSet)
	if err != nil {
		return nil, nil, err
	}

	// 4.Encode
	encoded, err := encodeData(data, canonicalMap)
	if err != nil {
		return nil, nil, err
	}

	return encoded, canonicalSet, nil
}

func encodeData(data []byte, codes map[byte]string) ([]byte, error) {
	var bitBuffer uint8
	var bitCount uint8
	var encoded []byte

	for _, b := range data {
		code, ok := codes[b]
		if !ok {
			return nil, fmt.Errorf("symbol %d not in codes", b)
		}

		for _, bit := range code {
			if bit == '1' {
				bitBuffer |= 1 << (7 - bitCount)
			}
			bitCount++

			if bitCount == 8 {
				encoded = append(encoded, bitBuffer)
				bitBuffer = 0
				bitCount = 0
			}
		}
	}

	if bitCount > 0 {
		encoded = append(encoded, bitBuffer)
	}

	return encoded, nil
}

func buildHuffmanTree(data []byte) *Node {
	// 1. Count frequency
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}

	// 2. Build nodes
	var nodes []*Node
	for symbol, count := range freq {
		nodes = append(nodes, &Node{Symbol: symbol, Frequency: count})
	}

	// 3. Build tree
	for len(nodes) > 1 {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Frequency < nodes[j].Frequency
		})

		left := nodes[0]
		right := nodes[1]
		nodes = nodes[2:]

		parent := &Node{
			Left:      left,
			Right:     right,
			Frequency: left.Frequency + right.Frequency,
		}
		nodes = append(nodes, parent)
	}

	return nodes[0]
}

func generateCodes(root *Node, prefix string, codes map[byte]string) {
	if root == nil {
		return
	}

	if root.Left == nil && root.Right == nil {
		codes[root.Symbol] = prefix
		return
	}

	generateCodes(root.Left, prefix+"0", codes)
	generateCodes(root.Right, prefix+"1", codes)
}

func Decode(encoded []byte, codes []Code, originalSize int) ([]byte, error) {
	codeMap, err := CodesToMap(codes)
	if err != nil {
		return nil, err
	}
	root := rebuildHuffmanTree(codeMap)
	var decoded []byte
	current := root
	bitPos := 0
	totalBits := len(encoded) * 8

	for bitPos < totalBits {
		bytePos := bitPos / 8
		bitOffset := 7 - (bitPos % 8)
		bit := (encoded[bytePos] >> bitOffset) & 1

		if bit == 0 {
			current = current.Left
		} else {
			current = current.Right
		}

		if current == nil {
			return nil, fmt.Errorf("invalid huffman code at bit %d", bitPos)
		}

		if current.Left == nil && current.Right == nil {
			decoded = append(decoded, current.Symbol)
			current = root

			if len(decoded) == originalSize {
				break
			}
		}

		bitPos++
	}

	if len(decoded) != originalSize {
		return nil, fmt.Errorf("decoded size mismatch: got %d, expected %d", len(decoded), originalSize)
	}

	return decoded, nil
}

func rebuildHuffmanTree(codes map[byte]string) *Node {
	root := &Node{}

	for symbol, code := range codes {
		current := root
		for _, bit := range code {
			if bit == '0' {
				if current.Left == nil {
					current.Left = &Node{}
				}
				current = current.Left
			} else {
				if current.Right == nil {
					current.Right = &Node{}
				}
				current = current.Right
			}
		}
		current.Symbol = symbol
	}

	return root
}
