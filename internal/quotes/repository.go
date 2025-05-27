package quotes

import (
	"encoding/json"
	"fmt"
	"go_text_task/pkg/db"
	"log"
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

	for node := r.Database.DL.Head; node != nil; node = node.Next {
		dbIndex := node.Value
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

	var id uint64
	for {
		id = newRand.Uint64()%r.Database.ID + 1
		if _, ok := r.Database.IndexMap[id]; ok {
			break
		}
	}

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

	for node := r.Database.DL.Head; node != nil; node = node.Next {
		dbIndex := node.Value
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

func (r *Repository) Create(payload QuoteRequest) (Quote, error) {
	quote := Quote{
		Author:    payload.Author,
		QuoteText: payload.QuoteText,
	}
	quote.CreatedAt = time.Now().Format(layout)
	quote.UpdatedAt = time.Now().Format(layout)
	quote.ID = r.Database.GetID() + 1
	data, encodeErr := json.Marshal(quote)
	if encodeErr != nil {
		return Quote{}, encodeErr
	}
	createErr := r.Database.Create(data)
	if createErr != nil {
		return Quote{}, createErr
	}
	return quote, nil
}

func (r *Repository) Delete(id uint64) error {
	return r.Database.Delete(id)
}

func (r *Repository) Update(id uint64, payload QuoteRequest) (Quote, error) {

	quote, getErr := r.GetQuoteByID(id)
	if getErr != nil {
		return Quote{}, getErr
	}

	quote.Author = payload.Author
	quote.QuoteText = payload.QuoteText
	quote.UpdatedAt = time.Now().Format(layout)

	go func() {
		var data []byte
		var updateErr error
		data, updateErr = json.Marshal(quote)
		updateErr = r.Database.Update(id, data)
		if updateErr != nil {
			log.Println("updating error: ", updateErr)
		}
	}()

	return quote, nil

}
