package services

import (
	"bytes"
	"compress/zlib"

	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/types"
)

type CanvasService struct {
	queries *data.Queries
}

func (s *CanvasService) compressPixelData(pixelData []types.Pixel) ([]byte, error) {
	var rawData bytes.Buffer
	for _, pixel := range pixelData {
		rawData.WriteByte(pixel.R)
		rawData.WriteByte(pixel.G)
		rawData.WriteByte(pixel.B)
		rawData.WriteByte(pixel.A)
	}

	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	defer zw.Close()
	if _, err := zw.Write(rawData.Bytes()); err != nil {
		return []byte{}, err
	}

	return compressed.Bytes(), nil
}

func (s *CanvasService) loadCanvas(compressed []byte) ([]types.Pixel, error) {
	var pixelArr []types.Pixel
	var decompressed bytes.Buffer

	zr, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return pixelArr, err
	}

	_, err = decompressed.ReadFrom(zr)
	if err != nil {
		return pixelArr, err
	}
	zr.Close()

	pixelBytes := decompressed.Bytes()
	pixelArr = make([]types.Pixel, len(pixelBytes)/4)
	for i := range len(pixelArr) {
		pixelArr[i] = types.Pixel{
			R: pixelBytes[i*4],
			G: pixelBytes[i*4+1],
			B: pixelBytes[i*4+2],
			A: pixelBytes[i*4+3],
		}
	}

	return pixelArr, nil
}

func (s *CanvasService) CreateCanvas(title, description string, width, height uint16, userID string) (types.Canvas, error) {
	pixelArr := make([]types.Pixel, width*height)
	pixelBytes, err := s.compressPixelData(pixelArr)
	if err != nil {
		return types.Canvas{}, err
	}

	canvas, err := s.queries.Canvas.CreateCanvas(title, description, userID, width, height, pixelBytes)
	if err != nil {
		return canvas, err
	}

	return canvas, nil
}
