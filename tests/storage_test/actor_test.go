package storage_test

import (
	"filmoteka/internal/actor"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestStorageAddActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	actorRepo := actor.NewActorRepository(db)

	testActor := &actor.Actor{Name: "Misha", Gender: "man", BirthDate: "1990-01-01"}

	//ok query
	mock.ExpectExec("INSERT INTO actor").
		WithArgs(testActor.Name, testActor.Gender, testActor.BirthDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = actorRepo.Add(testActor)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//query error
	mock.ExpectExec("INSERT INTO actor").
		WithArgs(testActor.Name, testActor.Gender, testActor.BirthDate).
		WillReturnError(fmt.Errorf("bad query"))

	err = actorRepo.Add(testActor)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestStorageGetActorId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	actorRepo := actor.NewActorRepository(db)

	testActor := &actor.Actor{
		Name:      "John Doe",
		Gender:    "man",
		BirthDate: "1990-01-01",
	}

	//Good query
	mock.ExpectQuery("SELECT id FROM actor").
		WithArgs(testActor.Name, testActor.Gender, testActor.BirthDate).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Возвращаем идентификатор актера

	actorID, err := actorRepo.GetActorId(testActor)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if actorID != 1 {
		t.Errorf("expected actorID 1, got %d", actorID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//query error
	mock.ExpectQuery("SELECT id FROM actor").
		WithArgs(testActor.Name, testActor.Gender, testActor.BirthDate).
		WillReturnError(fmt.Errorf("bad query"))

	actorID, err = actorRepo.GetActorId(testActor)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if actorID == 1 {
		t.Errorf("expected actorID 1, got %d", actorID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

}

func TestStorageUpdateActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	actorRepo := actor.NewActorRepository(db)

	actorID := int64(1)
	newActor := &actor.Actor{
		Name:      "John Doe",
		Gender:    "man",
		BirthDate: "1990-01-01",
	}

	mock.ExpectQuery("UPDATE actor").
		WithArgs(newActor.Name, newActor.Gender, newActor.BirthDate, actorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = actorRepo.Update(actorID, newActor)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//query error
	mock.ExpectQuery("UPDATE actor").
		WithArgs(newActor.Name, newActor.Gender, newActor.BirthDate, actorID).
		WillReturnError(fmt.Errorf("bad query"))

	err = actorRepo.Update(actorID, newActor)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

}

func TestStorageDeleteActor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	actorRepo := actor.NewActorRepository(db)

	actorID := int64(1)

	mock.ExpectQuery("DELETE FROM actor").
		WithArgs(actorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = actorRepo.Delete(actorID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//query error
	mock.ExpectQuery("DELETE FROM actor").
		WithArgs(actorID).
		WillReturnError(fmt.Errorf("bad query"))

	err = actorRepo.Delete(actorID)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
