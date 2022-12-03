package Schemas

type Quote struct {
	UUID   string `json:"UUID"`
	ID     int    `json:"ID"`
	Text   string `json:"text"`
	Author Author `json:"author"`
}
