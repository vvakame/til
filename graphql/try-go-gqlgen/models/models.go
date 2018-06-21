package models

type Todo struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID string `json:"-"`
}

type UserImpl struct {
	ID   string
	Name string
}
