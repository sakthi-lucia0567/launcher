package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("mad max")
	// Maximum upload size of 10 MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the uploads directory if it doesn't exist
	err = os.MkdirAll("assets/images", os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("assets/images/%s", handler.Filename))
	if err != nil {
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, handler.Filename)
}
