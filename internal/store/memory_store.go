package store

import (
	"errors"
	"sync"

	"github.com/leonardohenrique/pokemon-api/internal/models"
)

type PokemonStore struct {
	mu       sync.RWMutex
	pokemons map[int]models.Pokemon
	nextID   int
}

func NewPokemonStore() *PokemonStore {
	return &PokemonStore{
		pokemons: make(map[int]models.Pokemon),
		nextID:   1,
	}
}

func (s *PokemonStore) GetAll() []models.Pokemon {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Pokemon, 0, len(s.pokemons))
	for _, p := range s.pokemons {
		result = append(result, p)
	}
	return result
}

func (s *PokemonStore) GetByID(id int) (models.Pokemon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.pokemons[id]
	if !ok {
		return models.Pokemon{}, errors.New("pokemon não encontrado")
	}
	return p, nil
}

func (s *PokemonStore) Create(p models.Pokemon) models.Pokemon {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextID
	s.pokemons[p.ID] = p
	s.nextID++
	return p
}

func (s *PokemonStore) Update(id int, p models.Pokemon) (models.Pokemon, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.pokemons[id]; !ok {
		return models.Pokemon{}, errors.New("pokemon não encontrado")
	}
	p.ID = id
	s.pokemons[id] = p
	return p, nil
}

func (s *PokemonStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.pokemons[id]; !ok {
		return errors.New("pokemon não encontrado")
	}
	delete(s.pokemons, id)
	return nil
}
