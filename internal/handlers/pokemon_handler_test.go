package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leonardohenrique/pokemon-api/internal/models"
	"github.com/leonardohenrique/pokemon-api/internal/store"
)

type mockStore struct {
	pokemons []models.Pokemon
	total    int
	err      error
	lastTypes []string
}

func (m *mockStore) GetAll(ctx context.Context, page int, limit int, types []string) ([]models.Pokemon, int, error) {
	m.lastTypes = types
	return m.pokemons, m.total, m.err
}

func TestPokemonHandler_List_FilterTypes(t *testing.T) {
	s := &mockStore{
		pokemons: []models.Pokemon{
			{ID: 1, Name: "Pikachu", Type: []string{"electric"}, Level: 5, HP: 30},
			{ID: 2, Name: "Charmander", Type: []string{"fire"}, Level: 6, HP: 28},
		},
		total: 2,
		err:   nil,
	}
	h := &PokemonHandler{Store: s}

	req := httptest.NewRequest(http.MethodGet, "/pokemons?type=fire,electric", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Confirma que os tipos foram passados para o store
	if len(s.lastTypes) != 2 {
		t.Errorf("expected 2 filter types, got %d", len(s.lastTypes))
	}
	if s.lastTypes[0] != "fire" || s.lastTypes[1] != "electric" {
		t.Errorf("unexpected filter types: %v", s.lastTypes)
	}
}

func TestPokemonHandler_List_NoFilterTypes(t *testing.T) {
	s := &mockStore{
		pokemons: []models.Pokemon{},
		total:    0,
		err:      nil,
	}
	h := &PokemonHandler{Store: s}

	req := httptest.NewRequest(http.MethodGet, "/pokemons", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Confirma que nenhum tipo foi passado para o store
	if len(s.lastTypes) != 0 {
		t.Errorf("expected 0 filter types, got %d", len(s.lastTypes))
	}
}

func TestPokemonHandler_List_StoreError(t *testing.T) {
	s := &mockStore{
		pokemons: nil,
		total:    0,
		err:      errors.New("db error"),
	}
	h := &PokemonHandler{Store: s}

	req := httptest.NewRequest(http.MethodGet, "/pokemons?type=fire", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}
