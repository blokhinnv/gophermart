package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/asaskevich/govalidator"
)

func ReadBodyWithBodyReader(
	r *http.Request,
	allowedContentType string,
	bodyReader func(bodyBytes []byte) (any, error),
) (any, error) {
	if r.Header.Get("Content-Type") != allowedContentType {
		return nil, fmt.Errorf(
			"%w: incorrect content type",
			ErrIncorrectContentType,
		)
	}
	// проверим содержимое
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBytes) == 0 {
		return nil, fmt.Errorf(
			"%w: incorrent body (error while reading)",
			ErrIncorrectRequest,
		)
	}
	// читаем тело
	body, err := bodyReader(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: incorrent body (error while reading): %v",
			ErrIncorrectRequest,
			err.Error(),
		)
	}
	// валидируем данные
	if validated, err := govalidator.ValidateStruct(body); err != nil || !validated {
		return nil, fmt.Errorf(
			"%w: incorrent body (error while validating)",
			ErrNotValid,
		)
	}
	return body, nil
}
