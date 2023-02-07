package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/blokhinnv/gophermart/internal/app/database"
)

type logRegRequestBody struct {
	Login    string `json:"login"    valid:"required"`
	Password string `json:"password" valid:"required"`
}

const logRegBodyContentType = "application/json"

type LogReg struct {
	db             *database.DatabaseService
	signingKey     []byte
	expireDuration time.Duration
}

// чтение тела запроса с проверкой корректности
func (h *LogReg) ReadBody(r *http.Request) (*logRegRequestBody, int, error) {
	// проверим content-type
	if r.Header.Get("Content-Type") != logRegBodyContentType {
		return nil, http.StatusBadRequest, fmt.Errorf(
			"%w: incorrect content type",
			ErrIncorrectRequest,
		)
	}
	// проверим содержимое
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf(
			"%w: incorrent body (error while reading)",
			ErrIncorrectRequest,
		)
	}
	body := logRegRequestBody{}
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf(
			"%w: incorrent body (error while unmarshaling)",
			ErrIncorrectRequest,
		)
	}
	if validated, err := govalidator.ValidateStruct(body); err != nil || !validated {
		return nil, http.StatusBadRequest, fmt.Errorf(
			"%w: incorrent body (error while validating)",
			ErrIncorrectRequest,
		)
	}
	return &body, http.StatusOK, nil
}
