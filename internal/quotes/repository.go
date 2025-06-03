package quotes

import (
	"encoding/json"
	"fmt"
	"go_text_task/pkg/db"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"time"
)

const Layout = "2006-01-02 15:04:05"

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
		note, readErr := r.Database.Read(dbIndex)
		if readErr != nil {
			return []Quote{}, fmt.Errorf("database reading error: %v", readErr)
		}

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

	note, readErr := r.Database.Read(id)
	if readErr != nil {
		return Quote{}, fmt.Errorf("database reading error: %v", readErr)
	}

	var quote Quote
	decodeErr := json.Unmarshal(note, &quote)
	if decodeErr != nil {

		return Quote{}, fmt.Errorf("unmarshalling quote error: %v", decodeErr)
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
		note, readErr := r.Database.Read(dbIndex)
		if readErr != nil {
			return []Quote{}, fmt.Errorf("database reading error: %v", readErr)
		}

		var quote Quote
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

func (r *Repository) Create(payload *QuoteRequest) (Quote, error) {
	quote := Quote{
		Author:    payload.Author,
		QuoteText: payload.QuoteText,
	}
	quote.CreatedAt = time.Now().Format(Layout)
	quote.UpdatedAt = time.Now().Format(Layout)
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

func (r *Repository) Update(id uint64, payload *QuoteRequest) (Quote, error) {
	quote, getErr := r.GetQuoteByID(id)
	if getErr != nil {
		return Quote{}, getErr
	}

	quote.Author = payload.Author
	quote.QuoteText = payload.QuoteText
	quote.UpdatedAt = time.Now().Format(Layout)

	errGroup := errgroup.Group{}
	errGroup.Go(
		func() error {
			data, encodeErr := json.Marshal(quote)
			if encodeErr != nil {
				return encodeErr
			}
			return r.Database.Update(id, data)
		},
	)

	updateErr := errGroup.Wait()
	if updateErr != nil {
		return Quote{}, fmt.Errorf("updating error: %v", updateErr)
	}
	return quote, nil

}

func (r *Repository) Delete(id uint64) error {
	return r.Database.Delete(id)
}
