package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/cszczepaniak/go-htmx/internal/http/router"
	psql "github.com/cszczepaniak/go-htmx/internal/persistence/sql"
	"github.com/cszczepaniak/go-htmx/internal/sql"
)

//go:embed web/dist/*
var webAssets embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sql.NewFileDB("./data.db")
	if err != nil {
		panic(err)
	}

	p := psql.NewPersistence(db)
	err = p.Init(ctx)
	if err != nil {
		panic(err)
	}

	h := router.Setup(webAssets, p)

	f, err := os.OpenFile("./log/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("serving traffic!")

	err = http.ListenAndServe(":8080", h)
	if err != nil {
		panic(err)
	}

	// http.Handle("GET /events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	hdr := w.Header()
	// 	hdr.Set("Content-Type", "text/event-stream")
	// 	hdr.Set("Cache-Control", "no-cache")
	// 	hdr.Set("Connection", "keep-alive")
	//
	// 	t := time.NewTicker(time.Second)
	// 	for {
	// 		select {
	// 		case <-r.Context().Done():
	// 			fmt.Println("browser connection terminated")
	// 			return
	// 		case <-t.C:
	// 			// rand.Shuffle(len(items), func(i, j int) {
	// 			// 	items[i], items[j] = items[j], items[i]
	// 			// })
	// 			// fmt.Fprintln(w, "event: update")
	// 			// fmt.Fprint(w, "data: ")
	// 			// list(items).Render(context.Background(), w)
	// 			// fmt.Fprint(w, "\n\n")
	// 		}
	// 	}
	// }))
	//
	// http.ListenAndServe(":8080", nil)
}
