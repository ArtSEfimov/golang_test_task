package main

import (
	"fmt"
	"go_text_task/internal/quotes"
	"go_text_task/pkg/db"
	"net/http"
)

func startApp() {
	quotesMux := http.NewServeMux()

	quotesServer := &http.Server{
		Addr:    ":8080",
		Handler: quotesMux,
	}

	// DB init
	quotesDB := db.NewManager()

	// Quotes repository init
	quotesRepository := quotes.NewRepository(quotesDB)

	// Quotes handler init
	quotes.NewHandler(quotesMux, quotesRepository)

	fmt.Println("Quotes server listening at localhost:8080...")
	if err := quotesServer.ListenAndServe(); err != nil {
		fmt.Printf("quotes server err: %v\n", err)
	}
}

func main() {
	startApp()
}
