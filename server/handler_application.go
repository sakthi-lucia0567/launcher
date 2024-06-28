package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	internal "github.com/sakthi-lucia0567/launcher/internal/database"
)

func launcherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("hey html")
		http.ServeFile(w, r, "index.html")
	} else if r.Method == "POST" {

		type parameters struct {
			Name        string `json:"name"`
			Application string `json:"application"`
			Url         string `json:"url"`
		}

		decoder := json.NewDecoder(r.Body)

		params := parameters{}

		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}
		launchApplication(params.Application, params.Url)
		fmt.Fprintf(w, "Application launched: %s %s", params.Application, params.Url)
	}
}

func launchApplication(application, parameters string) {
	cmd := exec.Command(application, parameters)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error launching application:", err)
	}
}

func (apiCfg *apiConfig) handleCreateApplication(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string      `json:"name"`
		Path string      `json:"path"`
		Icon pgtype.Text `json:"icon"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	applicationUUID := uuid.New()
	generatedType := pgtype.UUID{Bytes: applicationUUID, Valid: true}

	application, err := apiCfg.DB.CreateApplication(r.Context(), internal.CreateApplicationParams{
		ID:        generatedType,
		Name:      params.Name,
		Path:      params.Path,
		Icon:      params.Icon,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create Application: %v", err))
		return
	}

	respondWithJSON(w, 201, (application))
}
