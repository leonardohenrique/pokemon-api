package main

import (
	"log"
	"net/http"
	"os"

	"github.com/leonardohenrique/pokemon-api/internal/handlers"
	"github.com/leonardohenrique/pokemon-api/internal/store"
)

func main() {
	pokemonStore := store.NewPokemonStore()
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
