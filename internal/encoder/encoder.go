package encoder

import (
	"bufio"
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

const blockSize = 50 * 1024 // 50KB

func EncodeFile(inputFile, outputFile string) error {
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

	if err := writeFileHeader(outFile); err != nil {
		return err
	}

	buffer := make([]byte, blockSize)
	for {
		n, err := inFile.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("read error: %w", err)
		}
		if n == 0 {
			break
		}

		compressed, meta, err := processBlock(buffer[:n])
		if err != nil {
			return fmt.Errorf("cannot process block: %w", err)
		}

		if err := writeBlock(outFile, compressed, meta); err != nil {
			return err
		}
	}

	return nil
}

func writeFileHeader(w io.Writer) error {
	header := []byte{'D', 'U', 'N', 'E'} // The Spice Must Flow
	_, err := w.Write(header)
	return err
}

// writeBlock use format [META_LEN:4][METADATA][DATA_LEN:4][COMPRESSED_DATA]
func writeBlock(w io.Writer, data []byte, meta *metadata.Metadata) error {
	metaBytes := meta.Serialize()

	buf := bufio.NewWriter(w)
	defer buf.Flush()

	// Write length of metadata (4 bytes)
	if err := binary.Write(buf, binary.BigEndian, uint32(len(metaBytes))); err != nil {
		return err
	}

	// Write metadata
	if _, err := buf.Write(metaBytes); err != nil {
		return err
	}

	// Write length of compressed data (4bytes)
	if err := binary.Write(buf, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}

	// Write compressed data
	if _, err := buf.Write(data); err != nil {
		return err
	}

	return nil
}

func processBlock(data []byte) ([]byte, *metadata.Metadata, error) {
	meta := &metadata.Metadata{
		BlockSize: uint32(len(data)),
	}

	// 1. Применяем BWT
	var bwtData []byte
	bwtData, meta.BWTPrimaryIndex = BWT.Encode(data)

	// 2. Применяем MTF
	mtfData := MTF.Encode(bwtData)

	// 3. Применяем RLE
	rleData := RLE.Encode(mtfData)
	meta.HuffmanSize = uint32(len(rleData))

	// 4. Применяем Huffman
	var err error
	var compressed []byte
	compressed, meta.HuffmanCodes, err = Huffman.Encode(rleData)
	if err != nil {
		return nil, nil, err
	}

	return compressed, meta, nil
}
