package api

type Book struct {
	ID             int          `json:"id"`
	Title          string       `json:"title"`
	Author         Author       `json:"author"`
	Series         Series       `json:"series"`
	Collections    []Collection `json:"collections"`
	SeriesNumber   int          `json:"seriesNumber"`
	CoverImagePath string       `json:"coverImagePath"`
}

// Rather than embedding, I separate the structs to make searching easier
// Example: I can search for all books in a series by searching for all books with the same Series.ID
// Additionally, I can find the names of all (unique) series easier

type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Series struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Collection is a group of books
// Example favorites, read, to-read
type Collection struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
