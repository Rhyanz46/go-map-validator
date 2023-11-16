package map_validator

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type FileValidation struct {
	AllowedTypes []FileType
	MaxMB        int64
}

type FileResult struct {
	Extension FileType
	Data      []byte
}

func getFileFromHttp(body io.ReadCloser, resWriter http.ResponseWriter, options FileValidation) (FileResult, int, error) {
	var extension FileType
	var res FileResult
	var valid bool
	body = http.MaxBytesReader(resWriter, body, options.MaxMB*1024*1024)

	data, err := io.ReadAll(body)
	if err != nil {
		if strings.Contains(err.Error(), "file is too large") {
			return res, http.StatusRequestEntityTooLarge, errors.New("file size is too large, maximal 1 MB")
		}
		return res, http.StatusInternalServerError, errors.New("internal server error")
	}
	for _, allowedType := range options.AllowedTypes {
		switch allowedType {
		case JPEG:
			valid, extension = isValidImageByte(data)
			if !valid {
				return res, http.StatusBadRequest, errors.New("file type is not valid image file")
			}
			if extension != JPEG {
				return res, http.StatusBadRequest, errors.New("image file must be jpeg format")
			}
		case PNG:
			valid, extension = isValidImageByte(data)
			if !valid {
				return res, http.StatusBadRequest, errors.New("file type is not valid image file")
			}
			if extension != PNG {
				return res, http.StatusBadRequest, errors.New("image file must be png format")
			}
		case IMAGE:
			valid, extension = isValidImageByte(data)
			if !valid {
				return res, http.StatusBadRequest, errors.New("file type is not valid image file")
			}
		default:
			return res, http.StatusBadRequest, errors.New("file type is not support yet to validate")
			// handle more type here
		}
	}
	return FileResult{
		Extension: extension,
		Data:      data,
	}, http.StatusOK, nil
}
