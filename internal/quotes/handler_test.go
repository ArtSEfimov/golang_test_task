package quotes_test

import (
	"bytes"
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
	"strings"
	"testing"
	"time"
)

const NOTES = 100_000

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

func TestGetAllHandler(t *testing.T) {
	setEnv(t)

	env := newTestFixture()
	testMux, testRepository := env.mux, env.repository

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	var w *httptest.ResponseRecorder
	w = httptest.NewRecorder()

	// empty DB test
	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.Bytes()

	var emptyQuotesList quotes.QuoteListResponse
	err := json.Unmarshal(body, &emptyQuotesList)
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

func TestGetByIDHandler(t *testing.T) {
	setEnv(t)

	env := newTestFixture()
	testMux, testRepository := env.mux, env.repository

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)

	// empty DB test
	testMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.Bytes()

	var emptyQuotesList quotes.QuoteListResponse
	err := json.Unmarshal(body, &emptyQuotesList)
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

func TestCreateHandler(t *testing.T) {
	setEnv(t)

	env := newTestFixture()
	testMux, testRepository := env.mux, env.repository

	for range NOTES {

		author := getRandomString(5)
		quoteText := getRandomString(10)
		requestQuote := quotes.QuoteRequest{
			Author:    author,
			QuoteText: quoteText,
		}
		bytesQuote, err := json.Marshal(&requestQuote)
		if err != nil {
			t.Fatal(err)
		}
		reader := bytes.NewReader(bytesQuote)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/quotes", reader)
		testMux.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("got %d, want %d", w.Code, http.StatusCreated)
		}
		body := w.Body.Bytes()
		var createdQuote quotes.Quote
		err = json.Unmarshal(body, &createdQuote)
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
			t.Fatalf("missmatch qoute QuoteText: got %s, want %s", createdQuote.QuoteText, dbQuote.QuoteText)
		}

	}
	defer func() {
		close(testRepository.Database.Tasks)
		<-testRepository.Database.Done
		_ = os.RemoveAll(filepath.Join(files.GetProjectRootDir(), "TestDatabase"))
	}()
}

func TestDeleteHandler(t *testing.T) {
	setEnv(t)

	env := newTestFixture()
	testMux, testRepository := env.mux, env.repository

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

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/quotes/%d", createdQuote.ID), nil)
		w := httptest.NewRecorder()
		testMux.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Fatalf("got %d, want %d", w.Code, http.StatusNoContent)
		}

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/quotes/%d", createdQuote.ID), nil)
		w = httptest.NewRecorder()
		testMux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("got %d, want %d", w.Code, http.StatusNotFound)
		}
		body := w.Body.Bytes()
		if strings.TrimSpace(string(body)) != fmt.Sprintf("quote %d not found", createdQuote.ID) {
			t.Fatalf("got %s, want %s", string(body), fmt.Sprintf("quote %d not found", createdQuote.ID))
		}

	}
	defer func() {
		close(testRepository.Database.Tasks)
		<-testRepository.Database.Done
		_ = os.RemoveAll(filepath.Join(files.GetProjectRootDir(), "TestDatabase"))
	}()
}

func TestUpdateHandler(t *testing.T) {
	setEnv(t)

	env := newTestFixture()
	testMux, testRepository := env.mux, env.repository

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

		author = getRandomString(5)
		quoteText = getRandomString(10)
		requestQuote := quotes.QuoteRequest{
			Author:    author,
			QuoteText: quoteText,
		}

		bytesQuote, err := json.Marshal(&requestQuote)
		if err != nil {
			t.Fatal(err)
		}
		reader := bytes.NewReader(bytesQuote)

		// for separation update time
		//time.Sleep(time.Second)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/quotes/%d", createdQuote.ID), reader)

		testMux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
		}

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/quotes/%d", createdQuote.ID), nil)
		w = httptest.NewRecorder()
		testMux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
		}

		body := w.Body.Bytes()
		var updatedQuote quotes.Quote
		err = json.Unmarshal(body, &updatedQuote)
		if err != nil {
			t.Fatal(err)
		}

		if updatedQuote.ID != createdQuote.ID {
			t.Fatalf("mismatch quote ID: got %d, want %d", updatedQuote.ID, createdQuote.ID)
		}
		if updatedQuote.Author != requestQuote.Author {
			t.Fatalf("mismatch quote Author: got %s, want %s", updatedQuote.Author, requestQuote.Author)
		}
		if updatedQuote.QuoteText != requestQuote.QuoteText {
			t.Fatalf("mismatch quote QuoteText: got %s, want %s", updatedQuote.QuoteText, requestQuote.QuoteText)
		}

		// You need to sleep for this test.
		//if updatedQuote.UpdatedAt == createdQuote.UpdatedAt {
		//	t.Fatalf("the creation time is equal to the update time, update time is %s, want %s", updatedQuote.UpdatedAt, createdQuote.UpdatedAt)
		//}

	}
	defer func() {
		close(testRepository.Database.Tasks)
		<-testRepository.Database.Done
		_ = os.RemoveAll(filepath.Join(files.GetProjectRootDir(), "TestDatabase"))
	}()
}
