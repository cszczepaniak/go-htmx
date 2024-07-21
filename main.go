package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func main() {
	num := 0
	component := shell()
	http.Handle("GET /", templ.Handler(component))
	http.Handle("POST /clicked", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		num++
		fmt.Fprintf(w, "clicked %d times", num)
	}))

	http.Handle("GET /events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := w.Header()
		hdr.Set("Content-Type", "text/event-stream")
		hdr.Set("Cache-Control", "no-cache")

		items := []string{"a", "b", "c", "d", "e", "f"}

		t := time.NewTicker(time.Second)
		for range t.C {
			rand.Shuffle(len(items), func(i, j int) {
				items[i], items[j] = items[j], items[i]
			})
			fmt.Fprintln(w, "event: update")
			fmt.Fprint(w, "data: ")
			list(items).Render(context.Background(), w)
			fmt.Fprint(w, "\n\n")
		}
	}))

	http.ListenAndServe(":8080", nil)
}
