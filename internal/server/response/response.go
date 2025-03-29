package response

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

type Map map[string]any

func EncodeDataToJSON(w http.ResponseWriter, r *http.Request, statuscode int, data any) error {
	if data == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statuscode)
		return nil
	}

	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statuscode)
	w.Write(j)
	return nil
}

func CreateHTTPErrorMessage(w http.ResponseWriter, r *http.Request, statuscode int, message string) error {
	msg := ErrorMessage{
		Message: message,
	}
	return EncodeDataToJSON(w, r, statuscode, msg)
}
