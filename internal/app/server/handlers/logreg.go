package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blokhinnv/gophermart/internal/app/database"
)

type logRegRequestBody struct {
	Login    string `json:"login"    valid:"required"`
	Password string `json:"password" valid:"required"`
}

const logRegBodyContentType = "application/json"

type LogReg struct {
	db             database.Service
	signingKey     []byte
	expireDuration time.Duration
}

// чтение тела запроса с проверкой корректности
func (h *LogReg) ReadBody(r *http.Request) (*logRegRequestBody, int, error) {
	bodyReader := func(bodyBytes []byte) (any, error) {
		body := logRegRequestBody{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, fmt.Errorf(
				"%w: incorrent body (error while unmarshaling)",
				ErrIncorrectRequest,
			)
		}
		return &body, nil
	}
	body, err := ReadBodyWithBodyReader(r, logRegBodyContentType, bodyReader)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if bodyTyped, ok := body.(*logRegRequestBody); ok {
		return bodyTyped, http.StatusOK, nil
	} else {
		return nil, http.StatusInternalServerError, nil
	}
}
