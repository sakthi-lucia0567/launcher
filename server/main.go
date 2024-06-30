package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	internal "github.com/sakthi-lucia0567/launcher/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *internal.Queries
}

func main() {

	// load the .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// conn, err := sql.Open("postgres", dbUrl)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	apiCfg := apiConfig{
		DB: internal.New(conn),
	}

	fs := http.FileServer(http.Dir("./assets"))

	router := chi.NewRouter()

	v1Router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://localhost:5173"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router.Get("/healthcheck", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// handling launcher
	v1Router.HandleFunc("/launcher", launcherHandler)
	v1Router.HandleFunc("/quit", quitHandler)

	v1Router.Post("/create_application", apiCfg.handleCreateApplication)
	v1Router.Get("/get_application", apiCfg.handleGetAllApplication)
	v1Router.Put("/update_application/{id}", apiCfg.handleUpdateApplication)
	v1Router.Delete("/delete_application/{id}", apiCfg.handleDeleteApplication)

	// File upload route
	v1Router.Post("/upload", handleFileUpload)

	router.Mount("/api/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server is running %v", port)
	error := srv.ListenAndServe()

	if err != nil {
		log.Fatal(error)
	}

}
