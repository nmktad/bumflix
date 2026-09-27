package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/nmktad/bumflix/internal/config"
	databse "github.com/nmktad/bumflix/internal/database"
	"github.com/nmktad/bumflix/internal/utils"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/movies/{prefix}", getMovie)
	http.ListenAndServe(":3000", r)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	prefixParam := chi.URLParam(r, "prefix")

	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error reading config: %v", err)
		return
	}

	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s?sslmode=disable", cfg.PGUser, cfg.PGPassword, cfg.PGDatabse))
	if err != nil {
		fmt.Printf("can't connect to db: %v", err)
		return 
	}

	defer conn.Close(ctx)

	queries := databse.New(conn)

	movie, err := queries.GetMovieByPrefix(ctx, prefixParam)
	if err != nil {
		w.WriteHeader(422)
		w.Write(fmt.Appendf(nil, "error fetching movie %s: %v", prefixParam, err))
		return
	}

	w.Write([]byte(utils.PlaybackURL(movie)))
}
