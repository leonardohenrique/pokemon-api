package store

import (
	"testing"

	"github.com/leonardohenrique/pokemon-api/internal/models"
)

func TestCreateAndGet(t *testing.T) {
	s := NewPokemonStore()
	created := s.Create(models.Pokemon{Name: "Charmander", Level: 5})

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("esperava encontrar, erro: %v", err)
	}
	if got.Name != "Charmander" {
		t.Errorf("esperava Charmander, recebeu %s", got.Name)
	}
}
