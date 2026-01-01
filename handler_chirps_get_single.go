package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGetSingle(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the chirpID string from the path
	chirpIDString := r.PathValue("chirpID")

	// 2. Parse the string into a UUID
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// 3. Retrieve the chirp from the database
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		// If SQLC returns no rows, send a 404
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
		return
	}

	// 4. Respond with 200 OK and the mapped Chirp
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
