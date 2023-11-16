package map_validator

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

func isValidImageByte(data []byte) (bool, FileType) {
	_, err := jpeg.Decode(bytes.NewReader(data))
	if err == nil {
		return true, JPEG
	}

	_, err = png.Decode(bytes.NewReader(data))
	if err == nil {
		return true, PNG
	}
	return false, ""
}
