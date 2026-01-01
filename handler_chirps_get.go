package main

import (
	"net/http"
	"sort"
	"time"

	"workspace/github.com/kozykoding/chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort") // Get the sort parameter

	var chirps []database.Chirp
	var err error

	if authorID != "" {
		uuidVal, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), uuidVal)
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	results := []Chirp{}
	for _, chirp := range chirps {
		results = append(results, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	// Sort the slice in-memory
	sort.Slice(results, func(i, j int) bool {
		if sortParam == "desc" {
			return results[i].CreatedAt.After(results[j].CreatedAt)
		}
		// Default to "asc"
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, results)
}
