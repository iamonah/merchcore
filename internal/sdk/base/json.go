package base

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IamOnah/storefronthq/internal/sdk/errs"

	"github.com/rs/zerolog/log"
)

// encoding
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	jsonData = append(jsonData, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

// decoding
func ReadJSON(r *http.Request, dst interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &dst)
	if err != nil {
		return err
	}
	return nil
}

// encode-errors
func WriteJSONError(w http.ResponseWriter, status errs.ErrCode, errors error) {
	errText := errs.CodeNames[status]
	code := errs.HTTPStatus[status]
	appError := errs.AppErr{
		Code:    code,
		Err:     errText,
		Details: errors.Error(),
	}
	err := WriteJSON(w, code, appError)
	if err != nil {
		log.Error().Err(err).Msg("writeJSONError")
		w.Header().Add("Connection", "close")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WriteJSONInternalError(w http.ResponseWriter, status errs.ErrCode) {
	errText := errs.CodeNames[status]
	code := errs.HTTPStatus[status]
	appError := errs.AppErr{
		Code: code,
		Err:  errText,
	}
	err := WriteJSON(w, code, appError)
	if err != nil {
		log.Error().Err(err).Msg("writeJSONError")
		w.Header().Add("Connection", "close")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
