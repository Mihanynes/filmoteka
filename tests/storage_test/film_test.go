package storage

import (
	"filmoteka/internal/actor"
	"filmoteka/internal/film"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestFilmRepositoryAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	filmRepo := film.NewFilmRepository(actor.NewActorRepository(db), db)

	film := &film.Film{
		Title:       "The Shawshank Redemption",
		Description: "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.",
		ReleaseDate: "1994-10-14",
		Rating:      9,
		Actors: []actor.Actor{
			{Name: "Tim Robbins", Gender: "man", BirthDate: "1958-10-16"},
			{Name: "Morgan Freeman", Gender: "man", BirthDate: "1937-06-01"},
		},
	}

	//good query
	mock.ExpectQuery("INSERT INTO film").
		WithArgs(film.Title, film.Description, film.ReleaseDate, film.Rating).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	for range film.Actors {
		mock.ExpectQuery("SELECT id FROM actor").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec("INSERT INTO film_actor").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	err = filmRepo.Add(film)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//query error
	mock.ExpectQuery("INSERT INTO film").
		WithArgs(film.Title, film.Description, film.ReleaseDate, film.Rating).
		WillReturnError(err)
	err = filmRepo.Add(film)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//query error

}

func TestFilmRepository_GetFilmId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	repo := film.NewFilmRepository(actor.NewActorRepository(db), db)
	// Mocked film data
	film := &film.Film{
		Title:       "TestFilm",
		ReleaseDate: "2023-01-01",
	}
	// good query
	rows := sqlmock.
		NewRows([]string{"id"}).
		AddRow(1)

	mock.
		ExpectQuery("SELECT id FROM film where title =").
		WithArgs(film.Title, film.ReleaseDate).
		WillReturnRows(rows)

	// Calling the method
	filmID, err := repo.GetFilmId(film)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	// Checking the result
	if filmID != 1 {
		t.Errorf("expected film ID 1, got %d", filmID)
	}

	// Checking if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// Query error
	mock.
		ExpectQuery("SELECT id FROM film where title =").
		WithArgs(film.Title, film.ReleaseDate).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetFilmId(film)
	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	// Checking if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// Row scan error
	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2) // Adding more rows to trigger row scan error

	mock.
		ExpectQuery("SELECT id FROM film where title =").
		WithArgs(film.Title, film.ReleaseDate).
		WillReturnRows(rows)

	_, err = repo.GetFilmId(film)
	if err != nil {
		t.Error("expected error, got nil")
		return
	}

	// Checking if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestFilmRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	repo := film.NewFilmRepository(actor.NewActorRepository(db), db)
	// Mocked film data
	filmID := int64(1)
	newFilm := &film.Film{
		Title:       "TestFilm",
		Description: "TestDescription",
		ReleaseDate: "01.01.2023",
		Rating:      9,
	}
	// good query
	mock.
		ExpectQuery("UPDATE film SET").
		WithArgs(newFilm.Title, newFilm.Description, newFilm.ReleaseDate, newFilm.Rating, filmID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Calling the method
	err = repo.Update(filmID, newFilm)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	// Checking if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// Query error
	mock.
		ExpectQuery("UPDATE film SET").
		WithArgs(newFilm.Title, newFilm.Description, newFilm.ReleaseDate, newFilm.Rating, filmID).
		WillReturnError(fmt.Errorf("db_error"))

	err = repo.Update(filmID, newFilm)
	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	// Checking if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestFilmRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	// Creating a new repository with mocked DB
	repo := film.NewFilmRepository(actor.NewActorRepository(db), db)

	// Mocking the transaction and queries
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM film_actor WHERE film_id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectExec("DELETE FROM film WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	// Calling the method
	err = repo.Delete(1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	// Checking if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// Query error
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM film_actor WHERE film_id = ?").
		WithArgs(2).
		WillReturnError(fmt.Errorf("db_error"))
	mock.ExpectRollback()

	// Calling the method
	err = repo.Delete(2)
	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	// Checking if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestFilmRepository_GetAllFilms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	// Creating a new repository with mocked DB
	repo := film.NewFilmRepository(actor.NewActorRepository(db), db)

	// Valid column test
	mock.ExpectPrepare("SELECT title, description, release_date, rating FROM film ORDER BY title")
	mock.ExpectQuery("SELECT title, description, release_date, rating FROM film ORDER BY title").
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "release_date", "rating"}))

	_, err = repo.GetAllFilms("title")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	// Invalid column test
	_, err = repo.GetAllFilms("invalid_column")
	if err == nil {
		t.Error("expected error, got nil for invalid column")
		return
	}
	if err.Error() != "film_repo.GetAllFilms: invalid column name: invalid_column" {
		t.Errorf("unexpected error message: %s", err)
		return
	}

	// Prepare error test
	mock.ExpectPrepare("SELECT title, description, release_date, rating FROM film ORDER BY title").
		WillReturnError(fmt.Errorf("prepare_error"))

	_, err = repo.GetAllFilms("title")
	if err == nil {
		t.Error("expected error, got nil for prepare error")
		return
	}
	if err.Error() != "film_repo.GetAllFilms: prepare_error" {
		t.Errorf("unexpected error message: %s", err)
		return
	}

	// Query execution error test
	mock.ExpectPrepare("SELECT title, description, release_date, rating FROM film ORDER BY title")
	mock.ExpectQuery("SELECT title, description, release_date, rating FROM film ORDER BY title").
		WillReturnError(fmt.Errorf("query_execution_error"))

	_, err = repo.GetAllFilms("title")
	if err == nil {
		t.Error("expected error, got nil for query execution error")
		return
	}
	if err.Error() != "film_repo.GetAllFilms: query_execution_error" {
		t.Errorf("unexpected error message: %s", err)
		return
	}

	// Row scan error test
	mock.ExpectPrepare("SELECT title, description, release_date, rating FROM film ORDER BY title")
	mock.ExpectQuery("SELECT title, description, release_date, rating FROM film ORDER BY title").
		WillReturnRows(sqlmock.NewRows([]string{"title"}).AddRow("film1"))

	_, err = repo.GetAllFilms("title")
	if err == nil {
		t.Error("expected error, got nil for row scan error")
		return
	}
	if err.Error() != "film_repo.GetAllFilms: sql: expected 1 destination arguments in Scan, not 4" {
		t.Errorf("unexpected error message: %s", err)
		return
	}
}

//func TestFilmRepository_FindFilms(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("can't create mock: %s", err)
//	}
//	defer db.Close()
//
//	// Creating a new repository with mocked DB
//	repo := film.NewFilmRepository(actor.NewActorRepository(db), db)
//
//	// Mocked input string
//	toFind := "Test"
//
//	// Mocked films
//	filmsByTitle := []film.Film{
//		{Title: "Test Film 1", Description: "Description 1", ReleaseDate: "2022-01-01", Rating: 8},
//		{Title: "Test Film 2", Description: "Description 2", ReleaseDate: "2022-01-02", Rating: 7},
//	}
//
//	// Mocking the database query for finding films by title
//	mock.ExpectQuery("SELECT f.title, f.description, f.release_date, f.rating FROM film f WHERE f.title LIKE ?").
//		WithArgs("%" + toFind + "%").
//		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "release_date", "rating"}).
//			AddRow(filmsByTitle[0].Title, filmsByTitle[0].Description, filmsByTitle[0].ReleaseDate, filmsByTitle[0].Rating).
//			AddRow(filmsByTitle[1].Title, filmsByTitle[1].Description, filmsByTitle[1].ReleaseDate, filmsByTitle[1].Rating))
//
//	// Calling the method
//	resultFilms, err := repo.FindFilms(toFind)
//	if err != nil {
//		t.Errorf("unexpected error: %s", err)
//		return
//	}
//
//	// Verifying the returned films
//	if !reflect.DeepEqual(resultFilms, filmsByTitle) {
//		t.Error("returned films do not match expected films")
//		return
//	}
//
//	// Checking if all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//
//	// Mocked films by actor
//	filmsByActor := []film.Film{
//		{Title: "Test Film 3", Description: "Description 3", ReleaseDate: "2022-01-03", Rating: 9},
//		{Title: "Test Film 4", Description: "Description 4", ReleaseDate: "2022-01-04", Rating: 8},
//	}
//
//	// Mocking the database query for finding films by actor
//	mock.ExpectQuery("SELECT f.title, f.description, f.release_date, f.rating FROM film f WHERE f.id IN").
//		WithArgs("%" + toFind + "%").
//		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "release_date", "rating"}).
//			AddRow(filmsByActor[0].Title, filmsByActor[0].Description, filmsByActor[0].ReleaseDate, filmsByActor[0].Rating).
//			AddRow(filmsByActor[1].Title, filmsByActor[1].Description, filmsByActor[1].ReleaseDate, filmsByActor[1].Rating))
//
//	// Calling the method with the same input string to test finding films by actor
//	resultFilms, err = repo.FindFilms(toFind)
//	if err != nil {
//		t.Errorf("unexpected error: %s", err)
//		return
//	}
//
//	// Verifying the returned films by actor
//	if !reflect.DeepEqual(resultFilms, filmsByActor) {
//		t.Error("returned films do not match expected films")
//		return
//	}
//
//	// Checking if all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//
//	// Mocking the scenario when no films are found
//	mock.ExpectQuery("SELECT f.title, f.description, f.release_date, f.rating FROM film f WHERE f.title LIKE ?").
//		WithArgs("%" + toFind + "%").
//		WillReturnRows(sqlmock.NewRows([]string{})) // Empty rows indicating no films found
//
//	mock.ExpectQuery("SELECT f.title, f.description, f.release_date, f.rating FROM film f WHERE f.id IN").
//		WithArgs("%" + toFind + "%").
//		WillReturnRows(sqlmock.NewRows([]string{})) // Empty rows indicating no films found
//
//	// Calling the method when no films are found
//	resultFilms, err = repo.FindFilms(toFind)
//	if err == nil {
//		t.Error("expected error, got nil")
//		return
//	}
//
//	// Verifying the error message
//	expectedErrorMessage := "film_repo.FindFilm: no films found"
//	if err.Error() != expectedErrorMessage {
//		t.Errorf("unexpected error message: got %s, expected %s", err.Error(), expectedErrorMessage)
//		return
//	}
//
//	// Checking if all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//}
