package store

import (
	"context"
	"errors"
	"fmt"
	"math"

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

func (s *PokemonStore) GetAll(ctx context.Context, page, limit int) ([]models.Pokemon, int, error) {
	offset := (page - 1) * limit

	// Total de registros, pra calcular quantas páginas existem
	var total int
	// Incluindo foto na consulta
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM pokemons").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar pokemons: %w", err)
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, name, type, level, hp, photo_url FROM pokemons
		 ORDER BY id
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar pokemons: %w", err)
	}
	defer rows.Close()

	result := make([]models.Pokemon, 0, limit)
	for rows.Next() {
		var p models.Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Level, &p.HP, &p.PhotoURL); err != nil {
			return nil, 0, fmt.Errorf("erro ao ler linha: %w", err)
		}
		result = append(result, p)
	}

	return result, total, nil
}

func TotalPages(totalItems, limit int) int {
	if limit == 0 {
		return 0
	}
	return int(math.Ceil(float64(totalItems) / float64(limit)))
}

func (s *PokemonStore) GetByID(ctx context.Context, id int) (models.Pokemon, error) {
	var p models.Pokemon
	err := s.db.QueryRow(ctx,
		"SELECT id, name, type, level, hp, photo_url FROM pokemons WHERE id = $1", id,
	).Scan(&p.ID, &p.Name, &p.Type, &p.Level, &p.HP, &p.PhotoURL)

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
		`INSERT INTO pokemons (name, type, level, hp, photo_url)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		p.Name, p.Type, p.Level, p.HP, p.PhotoURL,
	).Scan(&p.ID)

	if err != nil {
		return models.Pokemon{}, fmt.Errorf("failed to create pokemon: %w", err)
	}
	return p, nil
}

func (s *PokemonStore) Update(ctx context.Context, id int, p models.Pokemon) (models.Pokemon, error) {
	cmdTag, err := s.db.Exec(ctx,
		`UPDATE pokemons SET name = $1, type = $2, level = $3, hp = $4, photo_url = $5
		 WHERE id = $6`,
		p.Name, p.Type, p.Level, p.HP, p.PhotoURL, id,
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
