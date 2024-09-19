package models

type Book struct {
	ID             int          `json:"id"`
	Title          string       `json:"title"`
	Author         Author       `json:"author"`
	Series         Series       `json:"series"`
	Collections    []Collection `json:"collections"`
	SeriesNumber   int          `json:"seriesNumber"`
	CoverImagePath string       `json:"coverImagePath"`
}
