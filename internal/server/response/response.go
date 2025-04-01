package response

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	Error      *any   `json:"error"`
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

func CreateErrorResponse(w http.ResponseWriter, r *http.Request, statuscode int, errorObj any) error {
	resp := ErrorResponse{
		Status:     "Failed",
		StatusCode: strconv.Itoa(statuscode),
		Error:      &errorObj,
	}
	return EncodeDataToJSON(w, r, statuscode, resp)
}
