package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/leonardohenrique/pokemon-api/internal/models"
	"github.com/leonardohenrique/pokemon-api/internal/store"
	"github.com/leonardohenrique/pokemon-api/internal/validation"
)

type PokemonHandler struct {
	Store *store.PokemonStore
}

func NewPokemonHandler(s *store.PokemonStore) *PokemonHandler {
	return &PokemonHandler{Store: s}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// GET /pokemons
func (h *PokemonHandler) List(w http.ResponseWriter, r *http.Request) {
	pokemons, err := h.Store.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, pokemons)
}

// GET /pokemons/{id}
func (h *PokemonHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	pokemon, err := h.Store.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pokemon)
}

// POST /pokemons
func (h *PokemonHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Pokemon
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if validationErrors := validation.Validate(p); validationErrors != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error":  "invalid data",
			"fields": validationErrors,
		})
		return
	}

	created, err := h.Store.Create(r.Context(), p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// PUT /pokemons/{id}
func (h *PokemonHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var p models.Pokemon
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if validationErrors := validation.Validate(p); validationErrors != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error":  "invalid data",
			"fields": validationErrors,
		})
		return
	}

	updated, err := h.Store.Update(r.Context(), id, p)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// DELETE /pokemons/{id}
func (h *PokemonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
