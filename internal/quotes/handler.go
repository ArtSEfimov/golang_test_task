package quotes

import (
	"encoding/json"
	"fmt"
	"go_text_task/pkg/response"
	"net/http"
	"strconv"
)

type Handler struct {
	Repository *Repository
}

func NewHandler(router *http.ServeMux, repository *Repository) {
	handler := &Handler{
		Repository: repository,
	}
	router.HandleFunc("GET /quotes", handler.GetAll())
	router.HandleFunc("GET /quotes/random", handler.GetRandom())
	router.HandleFunc("GET /quotes/{id}", handler.GetByID())
	router.HandleFunc("POST /quotes", handler.Create())
	router.HandleFunc("DELETE /quotes/{id}", handler.Delete())
	router.HandleFunc("PUT /quotes/{id}", handler.Update())

}

func (handler *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var quotes []Quote
		var getErr error
		author := r.URL.Query().Get("author")
		if author == "" {
			quotes, getErr = handler.Repository.GetAll()
		} else {
			quotes, getErr = handler.Repository.GetByAuthor(author)
		}

		if getErr != nil {
			http.Error(w, getErr.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, QuoteListResponse{
			Quotes: quotes,
		}, http.StatusOK)
	}
}

func (handler *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
			return
		}

		quote, getErr := handler.Repository.GetQuoteByID(id)
		if getErr != nil {
			http.Error(w, getErr.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, quote, http.StatusOK)
	}
}

func (handler *Handler) GetRandom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		quote, getErr := handler.Repository.GetRandomQuote()
		if getErr != nil {
			http.Error(w, getErr.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, quote, http.StatusOK)
	}
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload QuoteRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&payload)
		if decodeErr != nil {
			e := fmt.Errorf("decoding error: %w", decodeErr)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		createdQuote, creationErr := handler.Repository.Create(&payload)
		if creationErr != nil {
			e := fmt.Errorf("creation error: %w", creationErr)
			http.Error(w, e.Error(), http.StatusBadGateway)
			return
		}

		response.Json(w, createdQuote, http.StatusCreated)
	}
}
func (handler *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
			return
		}

		deleteErr := handler.Repository.Delete(id)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, nil, http.StatusNoContent)
	}
}
func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
			return
		}

		var payload QuoteRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&payload)
		if decodeErr != nil {
			e := fmt.Errorf("decoding error: %w", decodeErr)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		updatedQuote, updateErr := handler.Repository.Update(id, &payload)
		if updateErr != nil {
			e := fmt.Errorf("updating error: %w", updateErr)
			http.Error(w, e.Error(), http.StatusBadGateway)
			return
		}

		response.Json(w, updatedQuote, http.StatusOK)
	}
}
