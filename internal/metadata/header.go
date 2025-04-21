package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Dadypool/archiver/pkg/algorithm/Huffman"
)

// Metadata is a block of metadata in front of compressed data
type Metadata struct {
	BlockSize uint32

	// BWT metadata
	BWTPrimaryIndex uint32

	// Huffman metadata
	HuffmanCodes []Huffman.Code
	HuffmanSize  uint32
}

func (m *Metadata) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, m.BlockSize)
	binary.Write(buf, binary.BigEndian, m.BWTPrimaryIndex)

	binary.Write(buf, binary.BigEndian, m.HuffmanSize)
	binary.Write(buf, binary.BigEndian, uint16(len(m.HuffmanCodes)))

	for _, code := range m.HuffmanCodes {
		buf.WriteByte(code.Symbol)
		buf.WriteByte(code.Length)
	}

	return buf.Bytes()
}

func DeserializeMetadata(data []byte) (*Metadata, error) {
	m := &Metadata{}
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.BigEndian, &m.BlockSize); err != nil {
		return nil, fmt.Errorf("failed to read BlockSize: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &m.BWTPrimaryIndex); err != nil {
		return nil, fmt.Errorf("failed to read BWTPrimaryIndex: %w", err)
	}
	if err := binary.Read(buf, binary.BigEndian, &m.HuffmanSize); err != nil {
		return nil, fmt.Errorf("failed to read HuffmanSize: %w", err)
	}

	var codesCount uint16
	if err := binary.Read(buf, binary.BigEndian, &codesCount); err != nil {
		return nil, fmt.Errorf("failed to read Huffman codes count: %w", err)
	}

	m.HuffmanCodes = make([]Huffman.Code, codesCount)
	for i := range m.HuffmanCodes {
		var err error
		if m.HuffmanCodes[i].Symbol, err = buf.ReadByte(); err != nil {
			return nil, fmt.Errorf("failed to read Huffman code symbol: %w", err)
		}
		if m.HuffmanCodes[i].Length, err = buf.ReadByte(); err != nil {
			return nil, fmt.Errorf("failed to read Huffman code length: %w", err)
		}
	}

	return m, nil
}
