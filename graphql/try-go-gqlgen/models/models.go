//go:generate dataloaden -keys string github.com/vvakame/til/graphql/try-go-gqlgen/models.UserImpl

package models

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

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

const UserLoaderKey = "userloader"

func DataloaderMiddleware(userMap map[string]UserImpl, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := UserImplLoader{
			maxBatch: 100,
			wait:     1 * time.Millisecond,
			fetch: func(ids []string) ([]*UserImpl, []error) {
				fmt.Println("UserImplLoader#fetch", ids)

				users := make([]*UserImpl, len(ids))
				for idx, id := range ids {
					user, ok := userMap[id]
					if ok {
						users[idx] = &user
					} else {
						users[idx] = nil
					}
				}

				return users, nil
			},
		}
		ctx := context.WithValue(r.Context(), UserLoaderKey, &userloader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
