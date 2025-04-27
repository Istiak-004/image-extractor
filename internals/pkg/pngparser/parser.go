package pngparser

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"strings"

	"github.com/istiak-004/image-extractor/internals/domain"
)

const (
	pngHeader = "\x89PNG\r\n\x1a\n"
)

type PNGParser interface {
	ParseBase64Image(base64image string) (*domain.ExtractedData, error)
}

type pngParser struct{}

func NewPNGParser() PNGParser {
	return &pngParser{}
}

func (p *pngParser) ParseBase64Image(base64Image string) (*domain.ExtractedData, error) {

	cleanBase64Data := cleanBase64Data(base64Image)

	imageData, err := base64.StdEncoding.DecodeString(cleanBase64Data)

	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}

	extractedData, err := parsePNGData(imageData)
	fmt.Println("------------extractedData--------------")
	fmt.Println(extractedData, err)
	fmt.Println("--------------------------")
	if err != nil {
		return nil, fmt.Errorf("failed to get image data : %w", err)
	}
	return extractedData, nil
}

func cleanBase64Data(base64Image string) string {
	if strings.HasPrefix(base64Image, "data:image/png;base64,") {
		return strings.SplitAfter(base64Image, ",")[1]
	}
	return base64Image
}

type pngChunk struct {
	Length uint32
	Type   [4]byte
	Data   []byte
	CRC    uint32
}

func readChunk(reader io.Reader) (*pngChunk, error) {
	var chunk pngChunk

	if err := binary.Read(reader, binary.BigEndian, &chunk.Length); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &chunk.Type); err != nil {
		return nil, err
	}

	chunk.Data = make([]byte, chunk.Length)
	if _, err := io.ReadFull(reader, chunk.Data); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &chunk.CRC); err != nil {
		return nil, err
	}

	if !verifyCRC(&chunk) {
		return nil, errors.New("PNG chunk CRC verification failed!")
	}

	return &chunk, nil
}

func verifyCRC(chunk *pngChunk) bool {
	crcWriter := crc32.NewIEEE()
	crcWriter.Write(chunk.Type[:])
	crcWriter.Write(chunk.Data)
	return crcWriter.Sum32() == chunk.CRC
}

func parsePNGData(data []byte) (*domain.ExtractedData, error) {
	reader := bytes.NewReader(data)

	// Verify PNG header
	header := make([]byte, 8)
	if _, err := io.ReadFull(reader, header); err != nil || string(header) != pngHeader {
		return nil, errors.New("invalid PNG file")
	}

	for {
		chunk, err := readChunk(reader)

		fmt.Println("####################")
		fmt.Println(string(chunk.Type[:]))
		fmt.Println("####################")

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading PNG chunk: %w", err)
		}

		if string(chunk.Type[:]) == "tEXt" {
			fmt.Println("----chunk.Data--->", string(chunk.Data))
			parts := bytes.SplitN(chunk.Data, []byte{0}, 2)
			fmt.Println("----parts--->", len(parts))

			if len(parts) == 2 {
				fmt.Println("======parts=====>>", string(parts[0]), string(parts[1]))
				if string(parts[0]) == "json" {
					var extracted domain.ExtractedData
					if err := json.Unmarshal(parts[1], &extracted); err != nil {
						// Continue searching other chunks if JSON is invalid
						continue
					}
					return &extracted, nil
				}
			}
		}
	}

	return nil, errors.New("no valid JSON data found in PNG chunks")
}
