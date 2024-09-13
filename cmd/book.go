package main

type Book struct {
	UUID           string `json:"id"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	Published      string `json:"published"`
	ISBN           string `json:"isbn"`
	Tags           []Tag  `json:"metadata"`
	CoverImagePath string `json:"coverImagePath"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
