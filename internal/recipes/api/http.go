package api

import (
	"fmt"
	"html"
	"net/http"
)

func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
