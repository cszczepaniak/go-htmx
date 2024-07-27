package httpwrap

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Renderer interface {
	// Render should render a result to the given writer. This could be an HTML template or any
	// other type of data being returned from the HTTP handler.
	Render(ctx context.Context, w io.Writer) error
}

type Request struct {
	Request  *http.Request
	Response http.ResponseWriter
}

type Handler[T Renderer] func(ctx context.Context, req Request) (T, error)

type Middleware[T Renderer] func(Handler[T]) Handler[T]

func Handle[T Renderer](
	m *http.ServeMux,
	pattern string,
	handler Handler[T],
	middleware ...Middleware[T],
) {
	m.Handle(pattern, WrapHandler(handler, middleware...))
}

func WrapHandler[T Renderer](
	handler Handler[T],
	mw ...Middleware[T],
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
			Response: w,
		}

		renderer, err := h(ctx, req)
		if err != nil {
			// Gotta set up logging!
			fmt.Println("error rendering", err)

			w.WriteHeader(statusFromError(err))
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				/// TODO Log!
			}

			return
		}

		err = renderer.Render(r.Context(), w)
		if err != nil {
			// Gotta set up logging!
		}
	}
}

func statusFromError(_ error) int {
	return http.StatusInternalServerError
}
