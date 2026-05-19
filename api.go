package main

import (
	"encoding/json"
	"net/http"
	"groupie-tracker/structs"
)

func getArtists() ([]structs.Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []structs.Artist

	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func getRelations() ([]structs.Relation, error) {

	url := "https://groupietrackers.herokuapp.com/api/relation"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relationResponse structs.RelationResponse

	err = json.NewDecoder(resp.Body).Decode(&relationResponse)
	if err != nil {
		return nil, err
	}

	return relationResponse.Index, nil
}

