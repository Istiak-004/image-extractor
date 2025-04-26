package pngparser

import (
	"encoding/base64"
	"errors"

	"github.com/istiak-004/image-extractor/internals/domain"
)

type PNGParser interface {
	ParseBase64Image(base64image string) (*domain.ExtractedData, error)
}

type pngParser struct{}

func NewPNGParser() PNGParser {
	return &pngParser{}
}

func (p *pngParser) ParseBase64Image(base64Image string) (*domain.ExtractedData, error) {
	_, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, errors.New("failed to decode base64 image")
	}

	return &domain.ExtractedData{
		Name:         "Jane Smith",
		Organization: "Beta Inc",
		Address:      "456 Elm Ave",
		Mobile:       "+1 555 5678",
	}, nil
}
