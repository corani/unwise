package web

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

type HighlightCategory string

const (
	HighlightCategoryBooks HighlightCategory = "books"
)

// Used by Moon+ Reader

type CreateHighlightRequest struct {
	Highlights []CreateHighlight `json:"highlights"`
}

type CreateHighlight struct {
	Title        string            `json:"title"`         // [book] book title
	Author       string            `json:"author"`        // [book] book author
	SourceURL    string            `json:"source_url"`    // [book] url of article / tweet / podcast
	Category     HighlightCategory `json:"category"`      // [book]
	Text         string            `json:"text"`          // [highlight] highlight text
	Note         string            `json:"note"`          // [highlight] note for highlight
	Chapter      string            `json:"chapter"`       // [highlight] chapter of the book
	Location     int               `json:"location"`      // [highlight] highlights location, used to order the highlights
	HighlightURL string            `json:"highlight_url"` // [highlight] unique url for the highlight (e.g. a concrete tweet)
}

type CreateHighlightResponse struct {
	ID                 int               `json:"id"`                  // [book] generated
	Title              string            `json:"title"`               // [book]
	Author             string            `json:"author"`              // [book]
	Category           HighlightCategory `json:"category"`            // [book]
	NumHighlights      int               `json:"num_highlights"`      // [book] calculated from highlights
	LastHighlightAt    string            `json:"last_highlight_at"`   // [book] calculated from highlights
	UpdatedAt          string            `json:"updated"`             // [book] generated
	SourceURL          string            `json:"source_url"`          // [book]
	ModifiedHighlights []int             `json:"modified_highlights"` // [highlight] generated
}

// Used by Obsidian Plugin

type ListHighlightsResponse struct {
	Results []ListHighlight `json:"results"`
}

type ListHighlight struct {
	ID        int    `json:"id"`      // generated
	BookID    int    `json:"book_id"` // foreign key to Book.ID
	Text      string `json:"text"`
	Note      string `json:"note"`
	Chapter   string `json:"chapter"`
	Location  int    `json:"location"`
	URL       string `json:"url"`
	UpdatedAt string `json:"updated"` // generated
}

type ListBooksResponse struct {
	Results []ListBook `json:"results"`
}

type ListBook struct {
	ID            int               `json:"id"` // generated
	Title         string            `json:"title"`
	Author        string            `json:"author"`
	NumHighlights int               `json:"num_highlights"` // calculated
	SourceURL     string            `json:"source_url"`
	Category      HighlightCategory `json:"category"`
	UpdatedAt     string            `json:"updated"` // generated
}
