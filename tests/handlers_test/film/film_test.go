package film

import (
	"bytes"
	"encoding/json"
	"filmoteka/internal/actor"
	"filmoteka/internal/film"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFilmHandler_AddFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	handler := &film.FilmHandler{
		FilmRepo: mockStorage,
	}

	testFilm := &film.Film{
		Title:       "Test Film",
		Description: "Test description",
		ReleaseDate: "20.03.2024",
		Rating:      8,
	}
	mockStorage.EXPECT().Add(testFilm).Return(nil)

	reqBody, err := json.Marshal(testFilm)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/films", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.AddFilm(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedResponse := "film added: {Test Film Test description 20.03.2024 8 []}"
	if body := strings.TrimSpace(w.Body.String()); body != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, body)
	}
}

func TestFilmHandler_UpdateFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)

	handler := &film.FilmHandler{
		FilmRepo: mockStorage,
	}

	oldFilm := film.Film{
		Title:       "Old Film",
		Description: "Old Description",
		ReleaseDate: "01.01.2022",
		Rating:      8,
		Actors: []actor.Actor{
			{Name: "Actor 1", BirthDate: "01.01.1990"},
			{Name: "Actor 2", BirthDate: "02.01.1990"},
		},
	}
	newFilm := film.Film{
		Title:       "New Film",
		Description: "New Description",
		ReleaseDate: "01.01.2023",
		Rating:      9,
		Actors: []actor.Actor{
			{Name: "Actor 3", BirthDate: "03.01.1990"},
			{Name: "Actor 4", BirthDate: "04.01.1990"},
		},
	}

	filmInfo := []film.Film{oldFilm, newFilm}
	filmJSON, err := json.Marshal(filmInfo)
	if err != nil {
		t.Fatalf("failed to marshal film data: %v", err)
	}

	mockStorage.EXPECT().GetFilmId(&oldFilm).Return(int64(1), nil)
	mockStorage.EXPECT().Update(int64(1), &newFilm).Return(nil)

	req, err := http.NewRequest("POST", "/films", bytes.NewReader(filmJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler.UpdateFilm(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expectedResponse := "film updated"
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, rr.Body.String())
	}
}

func TestFilmHandler_DeleteFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)

	handler := &film.FilmHandler{
		FilmRepo: mockStorage,
	}

	filmToDelete := film.Film{
		Title:       "Film to delete",
		Description: "Description to delete",
		ReleaseDate: "01.01.2022",
		Rating:      8,
		Actors: []actor.Actor{
			{Name: "Actor 1", BirthDate: "01.01.1990"},
			{Name: "Actor 2", BirthDate: "02.01.1990"},
		},
	}
	filmJSON, err := json.Marshal(filmToDelete)
	if err != nil {
		t.Fatalf("failed to marshal film data: %v", err)
	}

	mockStorage.EXPECT().GetFilmId(&filmToDelete).Return(int64(1), nil)
	mockStorage.EXPECT().Delete(int64(1)).Return(nil)

	req, err := http.NewRequest("POST", "/films", bytes.NewReader(filmJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler.DeleteFilm(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expectedResponse := "film deleted"
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, rr.Body.String())
	}
}

func TestFilmHandler_GetAllFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)

	handler := &film.FilmHandler{
		FilmRepo: mockStorage,
	}

	req, err := http.NewRequest("GET", "/films", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	mockStorage.EXPECT().GetAllFilms("title").Return([]film.Film{
		{Title: "Film 1", ReleaseDate: "01.01.2022", Rating: 8},
		{Title: "Film 2", ReleaseDate: "02.01.2022", Rating: 7},
	}, nil)

	rr := httptest.NewRecorder()

	handler.GetAllFilms(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expectedResponse := `[{"title":"Film 1","release_date":"01.01.2022","rating":8},{"title":"Film 2","release_date":"02.01.2022","rating":7}]`
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, rr.Body.String())
	}
}

func TestFilmHandler_FindFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)

	handler := &film.FilmHandler{
		FilmRepo: mockStorage,
	}

	req, err := http.NewRequest("GET", "/films?find=Test", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	mockStorage.EXPECT().FindFilms("Test").Return([]film.Film{
		{Title: "Test Film 1", ReleaseDate: "01.01.2022", Rating: 8},
		{Title: "Test Film 2", ReleaseDate: "02.01.2022", Rating: 7},
	}, nil)

	handler.FindFilms(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expectedResponse := `[{"title":"Test Film 1","release_date":"01.01.2022","rating":8},{"title":"Test Film 2","release_date":"02.01.2022","rating":7}]`
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, rr.Body.String())
	}
}
