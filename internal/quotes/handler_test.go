package quotes_test

import (
	"encoding/json"
	"fmt"
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

const NOTES = 5000

const CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var TestingRandSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

func setEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_DIR", "TestDatabase")
	t.Setenv("MAX_FILE_SEGMENT_SIZE", "100000")
}

type testFixture struct {
	mux        *http.ServeMux
	repository *quotes.Repository
}

func newTestFixture() *testFixture {
	testConfig := config.NewConfig()
	testMux := http.NewServeMux()
	testDB := db.NewManager(testConfig)
	testRepository := quotes.NewRepository(testDB)
	quotes.NewHandler(testMux, testRepository)
	return &testFixture{
		mux:        testMux,
		repository: testRepository,
	}
}

func TestGetAll(t *testing.T) {
	setEnv(t)

	fmt.Println("START 1 TEST")

	env := newTestFixture()
	testMux, testRepository := env.mux, env.repository

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	var w *httptest.ResponseRecorder
	w = httptest.NewRecorder()
	var err error

	// empty DB test
	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.Bytes()

	var emptyQuotesList quotes.QuoteListResponse
	err = json.Unmarshal(body, &emptyQuotesList)
	if err != nil {
		t.Fatal(err)
	}

	if emptyQuotesList.Quotes != nil {
		t.Fatalf(`QuotesList should be empty`)
	}

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

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	body = w.Body.Bytes()

	var filledQuotesList quotes.QuoteListResponse
	err = json.Unmarshal(body, &filledQuotesList)
	if err != nil {
		t.Fatal(err)
	}
	if len(filledQuotesList.Quotes) != NOTES {
		t.Fatalf("got %d quotes, want %d", len(filledQuotesList.Quotes), NOTES)
	}

	for _, quote := range filledQuotesList.Quotes {
		if _, ok := dataMap[quote.Author]; !ok {
			t.Fatalf(`Quote author not found in DB`)
		}

		if quote.QuoteText != dataMap[quote.Author] {
			t.Fatalf(`Missmatch author and quote in DB: got %s, want %s`, quote.QuoteText, dataMap[quote.Author])
		}

	}
	defer func() {
		close(testRepository.Database.Tasks)
		<-testRepository.Database.Done
		_ = os.RemoveAll(filepath.Join(files.GetProjectRootDir(), "TestDatabase"))
	}()

}

func getRandomString(length int) string {

	byteString := make([]byte, length)
	for i := 0; i < length; i++ {
		byteString[i] = CHARS[TestingRandSeed.Intn(len(CHARS))]
	}
	return string(byteString)

}

func TestGetByID(t *testing.T) {
	setEnv(t)
	fmt.Println("START 2 TEST")
	env := newTestFixture()
	testMux, testRepository := env.mux, env.repository

	w := httptest.NewRecorder()
	var err error

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)

	// empty DB test
	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.Bytes()

	var emptyQuotesList quotes.QuoteListResponse
	err = json.Unmarshal(body, &emptyQuotesList)
	if err != nil {
		t.Fatal(err)
	}

	if emptyQuotesList.Quotes != nil {
		t.Fatalf(`QuotesList should be empty`)
	}

	// add random notes one by one

	for range NOTES {

		author := getRandomString(5)
		quoteText := getRandomString(10)

		createdQuote, err := testRepository.Create(&quotes.QuoteRequest{
			Author:    author,
			QuoteText: quoteText,
		})
		if err != nil {
			t.Fatal(err)
		}

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/quotes/%d", createdQuote.ID), nil)

		w = httptest.NewRecorder()
		testMux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
		}

		body = w.Body.Bytes()

		var dbQuote quotes.Quote
		err = json.Unmarshal(body, &dbQuote)
		if err != nil {
			t.Fatal(err)
		}

		if createdQuote.ID != dbQuote.ID {
			t.Fatalf("missmatch quote ID: got %d, want %d", dbQuote.ID, createdQuote.ID)
		}

		if createdQuote.Author != dbQuote.Author {
			t.Fatalf("missmatch quote Author: got %s, want %s", createdQuote.Author, dbQuote.Author)
		}

		if createdQuote.QuoteText != dbQuote.QuoteText {
			t.Fatalf("missmatch quote QuoteText: got %s, want %s", createdQuote.QuoteText, dbQuote.QuoteText)
		}
	}
	defer func() {
		close(testRepository.Database.Tasks)
		<-testRepository.Database.Done
		_ = os.RemoveAll(filepath.Join(files.GetProjectRootDir(), "TestDatabase"))
	}()
}
