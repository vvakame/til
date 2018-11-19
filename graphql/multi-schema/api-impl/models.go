package api_impl

import (
	"fmt"
	"io"
	"strconv"
)

type Role string

const (
	RolePublic Role = "PUBLIC"
	RoleStaff  Role = "STAFF"
)

func (e Role) IsValid() bool {
	switch e {
	case RolePublic, RoleStaff:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

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
