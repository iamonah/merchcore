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
	// 204 means "No Content", so skip writing a body.
	if status == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err = w.Write(jsonData); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

// decoding
func ReadJSON(r *http.Request, dst interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, dst)
	if err != nil {
		return err
	}
	return nil
}

// encode-errors
func WriteJSONError(w http.ResponseWriter, errvalue error) {
	err, ok := errvalue.(*errs.AppErr)
	if ok {
		err := WriteJSON(w, err.Code, err)
		if err != nil {
			log.Error().Err(err).Msg("writejson")
			w.Header().Add("Connection", "close")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	WriteJSONInternalError(w)

}

func WriteJSONInternalError(w http.ResponseWriter) {
	appError := errs.AppErr{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	err := WriteJSON(w, http.StatusInternalServerError, appError)
	log.Error().Err(err).Msg("writejson")
	w.Header().Add("Connection", "close")
	w.WriteHeader(http.StatusInternalServerError)
}
