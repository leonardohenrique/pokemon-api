package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leonardohenrique/pokemon-api/internal/models"
	"github.com/leonardohenrique/pokemon-api/internal/store"
)

// Teste simples para verificar filtro por tipo na List
func TestPokemonHandler_List_FilterByType(t *testing.T) {
	store := store.NewPokemonStore(nil)
	h := NewPokemonHandler(store)

	// Stub para o método GetAll do store para testar filtro
	store.GetAll = func(ctxAny interface{}, page, limit int, pokeType string) ([]models.Pokemon, int, error) {
		if pokeType != "fire" {
			t.Errorf("esperava filtro 'fire', mas recebeu '%s'", pokeType)
		}
		return []models.Pokemon{{ID: 1, Name: "Charmander", Type: []string{"fire"}, Level: 10, HP: 39}}, 1, nil
	}

	req := httptest.NewRequest("GET", "/pokemons?type=fire", nil)
	w := httptest.NewRecorder()
	
	h.List(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("esperava status %d, recebeu %d", http.StatusOK, res.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Charmander") {
		t.Errorf("resposta não contém pokémon esperado")
	}
}
