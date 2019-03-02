package go_hello

import "net/http"

func Func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
