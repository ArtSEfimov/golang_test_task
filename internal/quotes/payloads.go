package quotes

type QuoteRequest struct {
	Author    string `json:"author"`
	QuoteText string `json:"quote"`
}

type QuoteListResponse struct {
	Quotes []Quote `json:"quotes"`
}
