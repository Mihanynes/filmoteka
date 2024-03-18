package storage

import (
	"filmoteka/internal/actor"
	"filmoteka/internal/film"
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
