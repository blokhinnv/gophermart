package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blokhinnv/gophermart/internal/app/database"
)

type Withdraw struct {
	db database.Service
}

type withdrawRequestBody struct {
	OrderID string  `valid:"orderID,required" json:"order"`
	Sum     float64 `valid:",required"        json:"sum"`
}

const withdrawContentType = "application/json"

func (h *Withdraw) ReadBody(r *http.Request) (*postOrderBody, int, error) {
	bodyReader := func(bodyBytes []byte) (any, error) {
		body := withdrawRequestBody{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, fmt.Errorf(
				"%w: incorrent body (error while unmarshaling)",
				ErrIncorrectRequest,
			)
		}
		return &body, nil
	}
	body, err := ReadBodyWithBodyReader(r, withdrawContentType, bodyReader)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if bodyTyped, ok := body.(*postOrderBody); ok {
		return bodyTyped, http.StatusOK, nil
	} else {
		return nil, http.StatusInternalServerError, nil
	}

}
