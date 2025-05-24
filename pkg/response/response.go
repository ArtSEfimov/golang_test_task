package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func Json(writer http.ResponseWriter, data interface{}, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	encodeErr := json.NewEncoder(writer).Encode(data)
	if encodeErr != nil {
		log.Println(encodeErr)
	}
}
