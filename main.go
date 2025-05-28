package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go_text_task/internal/quotes"
	"go_text_task/pkg/db"
	"go_text_task/pkg/db/config"
	"net/http"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func startApp() {
	quotesConfig := config.NewConfig()

	quotesMux := http.NewServeMux()

	quotesServer := &http.Server{
		Addr:    ":8080",
		Handler: quotesMux,
	}

	// DB init
	quotesDB := db.NewManager(quotesConfig)

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
