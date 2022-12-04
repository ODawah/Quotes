package Schemas

type Author struct {
	UUID string `json:"UUID"`
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

type QuoteList struct {
	Author Author  `json:"author"`
	Quotes []Quote `json:"quotes"`
}
