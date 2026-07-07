package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/leonardohenrique/pokemon-api/internal/models"
)

type PokemonStore struct {
	db *pgxpool.Pool
}

func NewPokemonStore(db *pgxpool.Pool) *PokemonStore {
	return &PokemonStore{db: db}
}

func (s *PokemonStore) GetAll(ctx context.Context) ([]models.Pokemon, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, name, type, level, hp FROM pokemons ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemons: %w", err)
	}
	defer rows.Close()

	var result []models.Pokemon
	for rows.Next() {
		var p models.Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Level, &p.HP); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, p)
	}
	return result, nil
}

func (s *PokemonStore) GetByID(ctx context.Context, id int) (models.Pokemon, error) {
	var p models.Pokemon
	err := s.db.QueryRow(ctx,
		"SELECT id, name, type, level, hp FROM pokemons WHERE id = $1", id,
	).Scan(&p.ID, &p.Name, &p.Type, &p.Level, &p.HP)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Pokemon{}, errors.New("pokemon not found")
	}
	if err != nil {
		return models.Pokemon{}, fmt.Errorf("failed to fetch pokemon: %w", err)
	}
	return p, nil
}

func (s *PokemonStore) Create(ctx context.Context, p models.Pokemon) (models.Pokemon, error) {
	err := s.db.QueryRow(ctx,
		`INSERT INTO pokemons (name, type, level, hp)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		p.Name, p.Type, p.Level, p.HP,
	).Scan(&p.ID)

	if err != nil {
		return models.Pokemon{}, fmt.Errorf("failed to create pokemon: %w", err)
	}
	return p, nil
}

func (s *PokemonStore) Update(ctx context.Context, id int, p models.Pokemon) (models.Pokemon, error) {
	cmdTag, err := s.db.Exec(ctx,
		`UPDATE pokemons SET name = $1, type = $2, level = $3, hp = $4
		 WHERE id = $5`,
		p.Name, p.Type, p.Level, p.HP, id,
	)
	if err != nil {
		return models.Pokemon{}, fmt.Errorf("failed to update pokemon: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return models.Pokemon{}, errors.New("pokemon not found")
	}

	p.ID = id
	return p, nil
}

func (s *PokemonStore) Delete(ctx context.Context, id int) error {
	cmdTag, err := s.db.Exec(ctx, "DELETE FROM pokemons WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete pokemon: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("pokemon not found")
	}
	return nil
}
