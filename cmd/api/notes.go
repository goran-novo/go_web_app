package main

import (
	"errors"
	"net/http"
	"strconv"

	"greenlight.alexedwards.net/internal/data"
	"greenlight.alexedwards.net/internal/validator"
)

func (app *application) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Text      string  `json:"text"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	note := &data.Note{
		UserID:    app.contextGetUser(r).ID,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Text:      input.Text,
	}

	v := validator.New()

	if data.ValidateNote(v, note); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Notes.Insert(note)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"note": note}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listNotesHandler handles the GET requests to fetch all notes for the logged-in user.
func (app *application) listNotesFromLoggedUserHandler(w http.ResponseWriter, r *http.Request) {
	// Assuming the user ID is stored in the context by the authentication middleware.
	user := app.contextGetUser(r)

	// Fetch notes from the database.
	notes, err := app.models.Notes.GetAllByUserID(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the notes data as a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"notes": notes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listAllNotes(w http.ResponseWriter, r *http.Request) {

	// Fetch notes from the database.
	notes, err := app.models.Notes.GetAllNotes()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the notes data as a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"notes": notes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) findNearbyForLoggedUserNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user from context, assuming middleware sets this
	user := app.contextGetUser(r)

	// Parse latitude and longitude from the query string
	queryValues := r.URL.Query()
	latitude, latErr := strconv.ParseFloat(queryValues.Get("latitude"), 64)
	longitude, lonErr := strconv.ParseFloat(queryValues.Get("longitude"), 64)

	if latErr != nil || lonErr != nil {
		app.badRequestResponse(w, r, errors.New("invalid latitude or longitude"))
		return
	}

	// Assume a default radius of 1000 meters (1 km)
	radius := 1000 // You could also make this a parameter if desired

	notes, err := app.models.Notes.FindNotesNearbyForUser(user.ID, latitude, longitude, radius)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"notes": notes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) findNearbyForSpecificCoordinatesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user from context, assuming middleware sets this
	user := app.contextGetUser(r)

	// Parse latitude and longitude from the query string
	queryValues := r.URL.Query()
	latitude, latErr := strconv.ParseFloat(queryValues.Get("latitude"), 64)
	longitude, lonErr := strconv.ParseFloat(queryValues.Get("longitude"), 64)

	if latErr != nil || lonErr != nil {
		app.badRequestResponse(w, r, errors.New("invalid latitude or longitude"))
		return
	}

	// Assume a default radius of 1000 meters (1 km)
	radius := 100 // You could also make this a parameter if desired

	notes, err := app.models.Notes.FindNotesNearbyForUser(user.ID, latitude, longitude, radius)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"notes": notes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
