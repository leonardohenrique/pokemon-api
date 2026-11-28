package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/leonardohenrique/pokemon-api/internal/models"
	"github.com/leonardohenrique/pokemon-api/internal/store"
	"github.com/leonardohenrique/pokemon-api/internal/validation"
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
	uploadDir    = "./uploads"
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

	page := parseQueryInt(r, "page", defaultPage)
	limit := parseQueryInt(r, "limit", defaultLimit)

	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	pokemons, total, err := h.Store.GetAll(r.Context(), page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	response := models.PaginatedResponse{
		Data:       pokemons,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: store.TotalPages(total, limit),
	}

	writeJSON(w, http.StatusOK, response)
}

func parseQueryInt(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
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

func saveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	fileName := fmt.Sprintf("pokemon_%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fileName, nil
}

// POST /pokemons
func (h *PokemonHandler) Create(w http.ResponseWriter, r *http.Request) {
	// aceitar json e/ou multipart
	contentType := r.Header.Get("Content-Type")
	var p models.Pokemon

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		jsonData := r.FormValue("pokemon")
		if err := json.Unmarshal([]byte(jsonData), &p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid pokemon json data")
			return
		}

		file, header, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()
			savedFileName, err := saveUploadedFile(file, header)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to save photo")
				return
			}
			p.PhotoURL = "/uploads/" + savedFileName
		}
	} else {
		// json normal
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
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
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		jsonData := r.FormValue("pokemon")
		if err := json.Unmarshal([]byte(jsonData), &p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid pokemon json data")
			return
		}

		file, header, err := r.FormFile("photo")
		if err == nil {
			defer file.Close()
			savedFileName, err := saveUploadedFile(file, header)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to save photo")
				return
			}
			p.PhotoURL = "/uploads/" + savedFileName
		}
	} else {
		// json normal
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
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
