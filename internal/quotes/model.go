package quotes

type Quote struct {
	ID        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Author    string `json:"author"`
	QuoteText string `json:"quote"`
}
