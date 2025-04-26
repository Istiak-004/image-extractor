package service

import (
	"strings"

	"github.com/istiak-004/image-extractor/internals/domain"
	"github.com/istiak-004/image-extractor/internals/pkg/pngparser"
)

type ExtractorService interface {
	ExtractFromBase64(base64Image string) (*domain.ExtractedData, error)
}

type extractorService struct {
	parser pngparser.PNGParser
}

func NewExtractorService() ExtractorService {
	return &extractorService{
		parser: pngparser.NewPNGParser(),
	}
}

func (s *extractorService) ExtractFromBase64(base64Image string) (*domain.ExtractedData, error) {
	// Remove data URL prefix if present
	if strings.HasPrefix(base64Image, "data:image/png;base64,") {
		base64Image = strings.SplitAfter(base64Image, ",")[1]
	}

	return s.parser.ParseBase64Image(base64Image)
}
