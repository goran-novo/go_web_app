package data

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"greenlight.alexedwards.net/internal/validator"
)

type Note struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int64     `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Text      string    `json:"text"`
}

type NoteModel struct {
	DB *sql.DB
}

func (m *NoteModel) Insert(note *Note) error {
	query := `
        INSERT INTO notes (user_id, latitude, longitude, text, location)
        VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($3, $2), 4326))
        RETURNING id, created_at`

	args := []any{note.UserID, note.Latitude, note.Longitude, note.Text}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&note.ID, &note.CreatedAt)
}

func (m *NoteModel) GetAllNotes() ([]*Note, error) {
	query := `
        SELECT id, created_at, latitude, longitude, text
        FROM notes`

	var notes []*Note

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.CreatedAt, &note.Latitude, &note.Longitude, &note.Text); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// GetAllByUserID retrieves all notes for a specific user
func (m *NoteModel) GetAllByUserID(userID int64) ([]*Note, error) {
	query := `
        SELECT id, created_at, latitude, longitude, text
        FROM notes
        WHERE user_id = $1`

	var notes []*Note

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.CreatedAt, &note.Latitude, &note.Longitude, &note.Text); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (m *NoteModel) FindNotesNearbyForUser(userID int64, latitude, longitude float64, radius int) ([]*Note, error) {
	query := `
        SELECT id, created_at, user_id, latitude, longitude, text
        FROM notes
        WHERE user_id != $1 AND ST_DWithin(location, ST_MakePoint($3, $2)::geography, $4)`

	var notes []*Note

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID, latitude, longitude, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.CreatedAt, &note.UserID, &note.Latitude, &note.Longitude, &note.Text); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (m *NoteModel) FindNotesNearbyByCoordinates(latitude, longitude float64, radius int) ([]*Note, error) {
	query := `
        SELECT id, created_at, user_id, latitude, longitude, text
        FROM notes
        WHERE ST_DWithin(location, ST_MakePoint($2, $1)::geography, $3)`

	var notes []*Note

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, latitude, longitude, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.CreatedAt, &note.UserID, &note.Latitude, &note.Longitude, &note.Text); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func ValidateNote(v *validator.Validator, note *Note) {
	v.Check(note.Text != "", "text", "must be provided")
	v.Check(len(note.Text) <= 500, "text", "must not be more than 500 characters")

	validateLatitude(v, note.Latitude)
	validateLongitude(v, note.Longitude)
}

// validateLatitude checks the latitude for being within -90 to 90 and for correct precision.
func validateLatitude(v *validator.Validator, latitude float64) {
	v.Check(latitude >= -90 && latitude <= 90, "latitude", "must be between -90 and 90")

	// Check precision, for example up to 6 decimal places
	latString := fmt.Sprintf("%.6f", latitude)
	latRegex := regexp.MustCompile(`^-?\d{1,3}\.\d{1,6}$`)
	v.Check(latRegex.MatchString(latString), "latitude", "must have up to 6 decimal places")
}

// validateLongitude checks the longitude for being within -180 to 180 and for correct precision.
func validateLongitude(v *validator.Validator, longitude float64) {
	v.Check(longitude >= -180 && longitude <= 180, "longitude", "must be between -180 and 180")

	// Check precision, for example up to 6 decimal places
	lonString := fmt.Sprintf("%.6f", longitude)
	lonRegex := regexp.MustCompile(`^-?\d{1,3}\.\d{1,6}$`)
	v.Check(lonRegex.MatchString(lonString), "longitude", "must have up to 6 decimal places")
}
