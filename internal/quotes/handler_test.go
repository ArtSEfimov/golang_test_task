package quotes_test

import (
	"encoding/json"
	"go_text_task/internal/quotes"
	"go_text_task/pkg/db"
	"go_text_task/pkg/db/config"
	"go_text_task/pkg/files"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const NOTES = 1000

const CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var TestingRandSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

func TestGetAll(t *testing.T) {
	t.Setenv("DATABASE_DIR", "TestDatabase")
	t.Setenv("MAX_FILE_SEGMENT_SIZE", "100000")
	testConfig := config.NewConfig()
	testMux := http.NewServeMux()
	testDB := db.NewManager(testConfig)
	testRepository := quotes.NewRepository(testDB)
	quotes.NewHandler(testMux, testRepository)

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	var w *httptest.ResponseRecorder
	w = httptest.NewRecorder()
	var err error

	// empty DB test
	testMux.ServeHTTP(w, req)
	bytes := w.Body.Bytes()
	var emptyQuotesList quotes.QuoteListResponse
	err = json.Unmarshal(bytes, &emptyQuotesList)
	if err != nil {
		t.Fatal(err)
	}
	if emptyQuotesList.Quotes != nil {
		t.Fatalf(`QuotesList should be empty`)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	_ = os.RemoveAll(filepath.Join(files.GetProjectRootDir(), "TestDatabase"))

	// add random notes
	dataMap := make(map[string]string, NOTES)

	for range NOTES {
		author := getRandomString(5)
		quoteText := getRandomString(10)
		dataMap[author] = quoteText

		_, err = testRepository.Create(&quotes.QuoteRequest{
			Author:    author,
			QuoteText: quoteText,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	w = httptest.NewRecorder()
	testMux.ServeHTTP(w, req)
	bytes = w.Body.Bytes()

	var filledQuotesList quotes.QuoteListResponse
	err = json.Unmarshal(bytes, &filledQuotesList)
	if err != nil {
		t.Fatal(err)
	}
	if len(filledQuotesList.Quotes) != NOTES {
		t.Fatalf("got %d quotes, want %d", len(filledQuotesList.Quotes), NOTES)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	for _, quote := range filledQuotesList.Quotes {
		if _, ok := dataMap[quote.Author]; !ok {
			t.Fatalf(`Quote author not found in DB`)
		}

		if quote.QuoteText != dataMap[quote.Author] {
			t.Fatalf(`Missmatch author and quote in DB`)
		}

	}
	_ = os.RemoveAll(filepath.Join(files.GetProjectRootDir(), "TestDatabase"))

}

func getRandomString(length int) string {

	byteString := make([]byte, length)
	for i := 0; i < length; i++ {
		time.Sleep(time.Millisecond)
		byteString[i] = CHARS[TestingRandSeed.Intn(len(CHARS))]
	}
	return string(byteString)

}
