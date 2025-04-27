package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"image"
	"io"
)

type PNGCreatorService interface {
	CreatePNGFromJson(imgData map[string]string) ([]byte, error)
}

type pngCreatorService struct{}

func NewPNGCreatorService() PNGCreatorService {
	return &pngCreatorService{}
}

func (c *pngCreatorService) CreatePNGFromJson(imgData map[string]string) ([]byte, error) {
	_ = image.NewRGBA(image.Rect(0, 0, 1, 1))

	jsonBytes, err := json.Marshal(imgData)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	buf.Write([]byte("\x89PNG\r\n\x1a\n"))

	// Write IHDR chunk
	writeChunk(buf, "IHDR", []byte{
		0x00, 0x00, 0x00, 0x01, // width
		0x00, 0x00, 0x00, 0x01, // height
		0x08, // bit depth
		0x02, // color type (RGB)
		0x00, // compression
		0x00, // filter
		0x00, // interlace
	})

	// Write JSON as tEXt chunk
	textData := append([]byte("json\x00"), jsonBytes...)
	writeChunk(buf, "tEXt", textData)

	// Write IEND chunk
	writeChunk(buf, "IEND", nil)

	return buf.Bytes(), nil
}

func writeChunk(w io.Writer, chunkType string, data []byte) {
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(data)))

	w.Write(length)
	w.Write([]byte(chunkType))
	w.Write(data)

	// Calculate CRC
	crc := crc32.NewIEEE()
	crc.Write([]byte(chunkType))
	crc.Write(data)

	crcBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(crcBytes, crc.Sum32())
	w.Write(crcBytes)
}
