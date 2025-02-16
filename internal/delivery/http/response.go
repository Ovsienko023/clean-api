package http

import (
	"encoding/json"
	"errors"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, code int, resp any) {
	if resp == nil {
		w.WriteHeader(code)
		return
	}

	response, err := json.Marshal(resp)
	if err != nil {
		errorContainer := ErrorResponse{}
		errorContainer.Done(w, errors.New("internal error"))
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(response)
}
