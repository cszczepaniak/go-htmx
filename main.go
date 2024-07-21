package main

import (
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := hello("Connor")
	http.Handle("GET /", templ.Handler(component))

	http.ListenAndServe(":8080", nil)
}
