package main

import (
	"groupie-tracker/structs"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func renderError(w http.ResponseWriter, code int, message string) {

	w.WriteHeader(code)

	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: message,
	}

	tmpl.Execute(w, data)
}

func home(w http.ResponseWriter, r *http.Request) {

	// HANDLE INVALID ROUTES
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "Page Not Found")
		return
	}

	artists, err := getArtists()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to fetch artists")
		return
	}

	// SEARCH INPUT
	search := strings.ToLower(r.URL.Query().Get("search"))

	// FILTER INPUTS
	creationFilter := r.URL.Query().Get("creation")
	albumFilter := r.URL.Query().Get("album")

	// FILTERED RESULTS
	var filteredArtists []structs.Artist

	for _, artist := range artists {

		match := true

		// SEARCH FILTER
		if search != "" {

			found := strings.Contains(strings.ToLower(artist.Name), search)

			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), search) {
					found = true
					break
				}
			}

			if !found {
				match = false
			}
		}

		// CREATION DATE FILTER
		if creationFilter != "" {

			creationYear, err := strconv.Atoi(creationFilter)
			if err == nil {

				if artist.CreationDate != creationYear {
					match = false
				}
			}
		}

		// FIRST ALBUM FILTER
		if albumFilter != "" {

			albumYear, err := strconv.Atoi(albumFilter)
			if err == nil {

				artistAlbumYear := artist.FirstAlbum[len(artist.FirstAlbum)-4:]

				albumInt, err := strconv.Atoi(artistAlbumYear)
				if err == nil {

					if albumInt != albumYear {
						match = false
					}
				}
			}
		}

		if match {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// LOAD TEMPLATE
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Template Error")
		return
	}

	// SEND DATA TO TEMPLATE
	tmpl.Execute(w, filteredArtists)
}

func artistHandler(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, http.StatusBadRequest, "Invalid Artist ID")
		return
	}

	artists, err := getArtists()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Error fetching artists")
		return
	}

	relations, err := getRelations()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Error fetching relations")
		return
	}

	var selectedArtist structs.Artist
	var selectedRelation structs.Relation

	// FIND ARTIST
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			break
		}
	}

	// HANDLE INVALID ARTIST
	if selectedArtist.ID == 0 {
		renderError(w, http.StatusNotFound, "Artist Not Found")
		return
	}

	// FIND RELATION
	for _, relation := range relations {
		if relation.ID == id {
			selectedRelation = relation
			break
		}
	}

	data := struct {
		Artist  structs.Artist
		Relation structs.Relation
	}{
		Artist:  selectedArtist,
		Relation: selectedRelation,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Template Error")
		return
	}

	tmpl.Execute(w, data)
}

func main() {

	// STATIC FILES
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// ROUTES
	http.HandleFunc("/", home)
	http.HandleFunc("/artist/", artistHandler)

	log.Println("Server is running on http://localhost:7070")

	http.ListenAndServe(":7070", nil)
}