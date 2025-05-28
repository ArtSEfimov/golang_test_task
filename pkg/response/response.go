package response

import (
	"encoding/json"
	"net/http"
)

func Json(writer http.ResponseWriter, data interface{}, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if statusCode == http.StatusNoContent {
		return
	}
	
	encodeErr := json.NewEncoder(writer).Encode(data)
	if encodeErr != nil {
		panic(encodeErr)
	}
}
