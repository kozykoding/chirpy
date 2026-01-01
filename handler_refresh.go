package main

import (
	"net/http"
	"time"

	"workspace/github.com/kozykoding/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Extract token from header ONLY
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	// Find user associated with the non-expired, non-revoked token
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Issue a NEW 1-hour access token (JWT)
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Extract token from header ONLY
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	// Revoke the token in the DB
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	// Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
