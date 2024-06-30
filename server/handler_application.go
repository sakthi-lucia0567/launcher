package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	internal "github.com/sakthi-lucia0567/launcher/internal/database"
)

var appMutex sync.Mutex
var runningApps = make(map[string]*exec.Cmd)

func launcherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("hey html")
		http.ServeFile(w, r, "index.html")
	} else if r.Method == "POST" {
		log.Println("matter")
		type parameters struct {
			Name        string `json:"name"`
			Application string `json:"path"`
			Url         string `json:"url"`
		}

		decoder := json.NewDecoder(r.Body)

		params := parameters{}

		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}
		// log.Printf("Application launched: %s %s %s", params.Name, params.Application, params.Url)
		// log.Println("urusai", params.Application)
		launchApplication(params.Name, params.Application, params.Url)
	}
}

func launchApplication(name, application, parameters string) {
	log.Printf("params with %s %s %s", name, application, parameters)
	cmd := exec.Command(application, parameters)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error launching application:", err)
	}
	appMutex.Lock()
	runningApps[name] = cmd
	appMutex.Unlock()
}

func quitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		type parameters struct {
			Name string `json:"name"`
		}

		decoder := json.NewDecoder(r.Body)

		params := parameters{}

		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		quitApplication(params.Name)
		fmt.Fprintf(w, "Application quit: %s", params.Name)
	} else {
		respondWithError(w, 405, "Method not allowed")
	}
}

func quitApplication(name string) {
	appMutex.Lock()
	defer appMutex.Unlock()

	cmd, exists := runningApps[name]
	if !exists {
		fmt.Println("No such application running:", name)
		return
	}

	err := cmd.Process.Kill()
	if err != nil {
		fmt.Println("Error quitting application:", err)
		return
	}

	delete(runningApps, name)
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

func (apiConfig *apiConfig) handleGetAllApplication(w http.ResponseWriter, r *http.Request) {

	applications, err := apiConfig.DB.ListApplication(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("couldn't Get the Applications %v", err))
		return
	}
	respondWithJSON(w, 201, databaseApplicationToApplication(applications))
}

func (apiCfg *apiConfig) handleUpdateApplication(w http.ResponseWriter, r *http.Request) {

	ID := chi.URLParam(r, "id")
	log.Printf("ID %s", ID)
	type parameters struct {
		Path string `json:"path"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	applicationUUID, err := uuid.Parse(ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid UUID format: %v", err))
		return
	}

	updatedType := pgtype.UUID{Bytes: applicationUUID, Valid: true}

	log.Printf("path %s", params.Path)

	application, err := apiCfg.DB.UpdateApplication(r.Context(), internal.UpdateApplicationParams{
		ID:        updatedType,
		Path:      params.Path,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update Application: %v", err))
		return
	}

	respondWithJSON(w, 200, application)
}

func (apiCfg *apiConfig) handleDeleteApplication(w http.ResponseWriter, r *http.Request) {

	ID := chi.URLParam(r, "id")

	applicationUUID, err := uuid.Parse(ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid UUID format: %v", err))
		return
	}

	applicationID := pgtype.UUID{Bytes: applicationUUID, Valid: true}

	err = apiCfg.DB.DeleteApplication(r.Context(), applicationID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete Application: %v", err))
		return
	}

	respondWithJSON(w, 200, map[string]string{"message": "Application deleted"})
}
