package main

import (
	"time"

	"github.com/google/uuid"
	internal "github.com/sakthi-lucia0567/launcher/internal/database"
)

type Application struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseApplicationsToApplications(dbApplication internal.Application) Application {
	return Application{
		ID:        dbApplication.ID.Bytes,
		Name:      dbApplication.Name,
		Path:      dbApplication.Path,
		Icon:      dbApplication.Icon.String,
		CreatedAt: dbApplication.CreatedAt.Time,
		UpdatedAt: dbApplication.UpdatedAt.Time,
	}
}

func databaseApplicationToApplication(dbApplication []internal.Application) []Application {
	applications := []Application{}
	for _, dbApplication := range dbApplication {
		applications = append(applications, databaseApplicationsToApplications(dbApplication))
	}
	return applications
}
