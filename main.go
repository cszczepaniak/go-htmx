package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/cszczepaniak/go-htmx/internal/player/model"
	"github.com/cszczepaniak/go-htmx/internal/player/persistence/sqlite"
	"github.com/cszczepaniak/go-htmx/internal/sql"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sql.NewMemoryDB()
	if err != nil {
		panic(err)
	}

	p := sqlite.NewSQLitePlayerPersistence(db)
	err = p.Init(ctx)
	if err != nil {
		panic(err)
	}

	component := shell()
	http.Handle("GET /", templ.Handler(component))
	http.Handle("POST /players", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("serving post")
		err := r.ParseForm()
		if err != nil {
			fmt.Println("error parsing form", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = p.InsertPlayer(ctx, r.FormValue("firstName"), r.FormValue("lastName"))
		if err != nil {
			fmt.Println("error inserting player", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var players []model.Player
		rows, err := db.QueryContext(ctx, `SELECT ID, FirstName, LastName FROM Players`)
		if err != nil {
			fmt.Println("error querying db", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var p model.Player
			err := rows.Scan(&p.ID, &p.FirstName, &p.LastName)
			if err != nil {
				fmt.Println("error scanning row", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			players = append(players, p)
		}

		if err := rows.Err(); err != nil {
			fmt.Println("error from rows", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("got here", players)

		list(players).Render(ctx, w)
	}))

	http.Handle("GET /events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := w.Header()
		hdr.Set("Content-Type", "text/event-stream")
		hdr.Set("Cache-Control", "no-cache")
		hdr.Set("Connection", "keep-alive")

		t := time.NewTicker(time.Second)
		for {
			select {
			case <-r.Context().Done():
				fmt.Println("browser connection terminated")
				return
			case <-t.C:
				// rand.Shuffle(len(items), func(i, j int) {
				// 	items[i], items[j] = items[j], items[i]
				// })
				// fmt.Fprintln(w, "event: update")
				// fmt.Fprint(w, "data: ")
				// list(items).Render(context.Background(), w)
				// fmt.Fprint(w, "\n\n")
			}
		}
	}))

	http.ListenAndServe(":8080", nil)
}
