package quotes

import (
	"encoding/json"
	"fmt"
	"go_text_task/pkg/db"
	"math/rand"
	"time"
)

const layout = "2006-01-02 15:04:05"

type Repository struct {
	Database *db.Manager
}

func NewRepository(database *db.Manager) *Repository {
	return &Repository{
		Database: database,
	}
}

func (r *Repository) GetAll() ([]Quote, error) {
	var quotes []Quote
	for dbIndex := range r.Database.IndexMap {
		var quote Quote
		note, _ := r.Database.Read(dbIndex)
		decodeErr := json.Unmarshal(note, &quote)
		if decodeErr != nil {
			return []Quote{}, fmt.Errorf("unmarshalling quote error: %v", decodeErr)
		}
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func (r *Repository) GetQuoteByID(id uint64) (Quote, error) {
	if _, ok := r.Database.IndexMap[id]; !ok {
		return Quote{}, fmt.Errorf("quote %d not found", id)
	}
	var quote Quote
	note, _ := r.Database.Read(id)
	decodeErr := json.Unmarshal(note, &quote)
	if decodeErr != nil {

		return quote, fmt.Errorf("unmarshalling quote error: %v", decodeErr)
	}

	return quote, nil
}

func (r *Repository) GetRandomQuote() (Quote, error) {
	rs := rand.NewSource(time.Now().UnixNano())
	newRand := rand.New(rs)
	id := newRand.Uint64()%r.Database.ID + 1

	var quote Quote
	note, _ := r.Database.Read(id)
	decodeErr := json.Unmarshal(note, &quote)
	if decodeErr != nil {

		return quote, fmt.Errorf("unmarshalling quote error: %v", decodeErr)
	}

	return quote, nil
}

func (r *Repository) GetByAuthor(author string) ([]Quote, error) {
	var quotes []Quote
	for dbIndex := range r.Database.IndexMap {
		var quote Quote
		note, _ := r.Database.Read(dbIndex)
		decodeErr := json.Unmarshal(note, &quote)
		if decodeErr != nil {
			return []Quote{}, fmt.Errorf("unmarshalling quote error: %v", decodeErr)
		}
		if author == quote.Author {
			quotes = append(quotes, quote)
		}
	}
	return quotes, nil
}

func (r *Repository) Create(quote Quote) (Quote, error) {
	quote.CreatedAt = time.Now().Format(layout)
	quote.UpdatedAt = time.Now().Format(layout)
	quote.ID = r.Database.GetID() + 1
	data, encodeErr := json.Marshal(quote)
	if encodeErr != nil {
		return Quote{}, encodeErr
	}
	creationErr := r.Database.Create(data)
	if creationErr != nil {
		return Quote{}, creationErr
	}
	return quote, nil
}

func (r *Repository) Delete(id uint64) error {
	return r.Database.Delete(id)
}
