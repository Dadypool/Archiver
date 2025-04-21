package decoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Dadypool/archiver/internal/metadata"
	"github.com/Dadypool/archiver/pkg/algorithm/BWT"
	"github.com/Dadypool/archiver/pkg/algorithm/Huffman"
	"github.com/Dadypool/archiver/pkg/algorithm/MTF"
	"github.com/Dadypool/archiver/pkg/algorithm/RLE"
)

func DecodeFile(inputFile, outputFile string) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer outFile.Close()

	if err := readFileHeader(inFile); err != nil {
		return err
	}

	bufferedOut := bufio.NewWriter(outFile)
	defer bufferedOut.Flush()

	for {
		originalData, err := readBlock(inFile)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("block read error: %w", err)
		}

		if _, err := bufferedOut.Write(originalData); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}

	return nil
}

func readFileHeader(r io.Reader) error {
	expectedHeader := []byte{'D', 'U', 'N', 'E'}
	actualHeader := make([]byte, 4)

	if _, err := io.ReadFull(r, actualHeader); err != nil {
		return fmt.Errorf("header read error: %w", err)
	}

	if !bytes.Equal(actualHeader, expectedHeader) {
		return fmt.Errorf("invalid file header")
	}

	return nil
}

func readBlock(r io.Reader) ([]byte, error) {
	var metaLen uint32
	if err := binary.Read(r, binary.BigEndian, &metaLen); err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("meta length read error: %w", err)
	}

	metaBytes := make([]byte, metaLen)
	if _, err := io.ReadFull(r, metaBytes); err != nil {
		return nil, fmt.Errorf("metadata read error: %w", err)
	}

	meta, err := metadata.DeserializeMetadata(metaBytes)
	if err != nil {
		return nil, fmt.Errorf("metadata deserialize error: %w", err)
	}

	var dataLen uint32
	if err := binary.Read(r, binary.BigEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("data length read error: %w", err)
	}

	compressedData := make([]byte, dataLen)
	if _, err := io.ReadFull(r, compressedData); err != nil {
		return nil, fmt.Errorf("compressed data read error: %w", err)
	}

	return processBlock(compressedData, meta)
}

func processBlock(data []byte, meta *metadata.Metadata) ([]byte, error) {
	// 1. Inverse Huffman
	huffmanDecoded, err := Huffman.Decode(data, meta.HuffmanCodes, int(meta.HuffmanSize))
	if err != nil {
		return nil, err
	}

	// 2. Inverse RLE
	rleDecoded := RLE.Decode(huffmanDecoded)

	// 3. Inverse MTF
	mtfDecoded := MTF.Decode(rleDecoded)

	// 4. Inverse BWT
	originalData, err := BWT.Decode(mtfDecoded, int(meta.BWTPrimaryIndex))
	if err != nil {
		return nil, err
	}

	if int(meta.BlockSize) < len(originalData) {
		originalData = originalData[:meta.BlockSize]
	}

	return originalData, nil
}
