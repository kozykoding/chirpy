package main

import (
	"database/sql"
	"net/http"

	"workspace/github.com/kozykoding/chirpy/internal/auth"
	"workspace/github.com/kozykoding/chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate via Access Token (JWT)
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// 2. Extract and Parse the Chirp ID from the path
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// 3. Fetch the chirp first to check ownership (and existence)
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
		return
	}

	// 4. Authorize: Only the author can delete
	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not the author of this chirp", nil)
		return
	}

	// 5. Delete the Chirp
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	// 6. Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
