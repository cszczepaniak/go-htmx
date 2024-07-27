package httpwrap

import (
	"context"
	"fmt"
	"net/http"
)

type Handler func(ctx context.Context, req Request) error

type Middleware func(Handler) Handler

func Handle(
	m *http.ServeMux,
	pattern string,
	handler Handler,
	middleware ...Middleware,
) {
	m.Handle(pattern, WrapHandler(handler, middleware...))
}

func WrapHandler(
	handler Handler,
	mw ...Middleware,
) http.HandlerFunc {
	// Wrap the middlewares by looping through backwards. By going backwards, the first middleware
	// in the list will be the first one called, so it retains the order passed to us.
	h := handler
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := Request{
			Request:  r,
			response: w,
		}

		err := h(ctx, req)
		if err != nil {
			// Gotta set up logging!
			fmt.Println("error rendering", err)

			status := StatusCodeForError(err)
			w.WriteHeader(status)

			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				// TODO Log!
			}
		}
	}
}
