package api_impl

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type Todo struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID string `json:"-"`
}

type NewUser struct {
	Name  string `json:"name"`
	Staff bool   `json:"staff"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Staff bool   `json:"staff"`
}
