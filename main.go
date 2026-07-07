package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/leonardohenrique/pokemon-api/internal/handlers"
	"github.com/leonardohenrique/pokemon-api/internal/store"
)

func main() {

	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pokemon:pokemon123@localhost:5432/pokemon_db"
	}

	db, err := store.NewPostgresPool(ctx, dsn)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	pokemonStore := store.NewPokemonStore(db)
	h := handlers.NewPokemonHandler(pokemonStore)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /pokemons", h.List)
	mux.HandleFunc("GET /pokemons/{id}", h.Get)
	mux.HandleFunc("POST /pokemons", h.Create)
	mux.HandleFunc("PUT /pokemons/{id}", h.Update)
	mux.HandleFunc("DELETE /pokemons/{id}", h.Delete)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("servidor rodando na porta %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
