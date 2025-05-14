package services

import (
	"bytes"
	"compress/zlib"
	"errors"

	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/types"
	"github.com/jackc/pgx/v5"
)

type CanvasService struct {
	queries *data.Queries
}

func (s *CanvasService) CompressPixelData(pixelData []types.Pixel) ([]byte, error) {
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

func (s *CanvasService) LoadCanvas(compressed []byte) ([]types.Pixel, error) {
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

func (s *CanvasService) UserHasAccess(canvasID, userID string) (types.UserAccess, types.Canvas, error) {
	canvas, err := s.queries.GetCanvas(canvasID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.UserAccess{}, canvas, types.ErrCanvasDoesNotExist
		}
		return types.UserAccess{}, canvas, err
	}

	userAccess, err := s.queries.GetUserAccess(canvas.ID, types.CanvasObject, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) && canvas.LinkAccessType == types.Restricted {
			return userAccess, canvas, types.ErrUserAccessDenied
		}

		return userAccess, canvas, err
	}

	return userAccess, canvas, nil
}
