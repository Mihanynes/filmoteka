package unit_test

import (
	"filmoteka/internal/actor"
	"filmoteka/internal/film"
	"filmoteka/pkg"
	"reflect"
	"testing"
)

func TestDateValidation_ValidDate(t *testing.T) {
	date := "01.01.2000"
	err := pkg.DateValidation(date)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDateValidation_InvalidDate(t *testing.T) {
	date := "32.13.2022" // Неверный день и месяц
	err := pkg.DateValidation(date)
	expectedErr := "32.13.2022 is not a valid date\n"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	date = "01.01.20" // Неполный год
	err = pkg.DateValidation(date)
	expectedErr = "01.01.20 is not a valid date\n"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	date = "01/01/2022" // Неверный разделитель
	err = pkg.DateValidation(date)
	expectedErr = "01/01/2022 is not a valid date\n"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}
}

func TestConvertMapToActorsListWithFilms(t *testing.T) {
	// Sample data
	data := map[actor.Actor][]film.Film{
		actor.Actor{Name: "Actor1", BirthDate: "12.03.1995"}: []film.Film{{Title: "Film1"}, {Title: "Film2"}},
		actor.Actor{Name: "Actor2", BirthDate: "10.05.1989"}: []film.Film{{Title: "Film3"}, {Title: "Film4"}},
	}

	// Expected result
	expected := []film.ActorListWithFilms{
		{
			ActorInfo: actor.Actor{Name: "Actor1", BirthDate: "12.03.1995"},
			Films:     []film.Film{{Title: "Film1"}, {Title: "Film2"}},
		},
		{
			ActorInfo: actor.Actor{Name: "Actor2", BirthDate: "10.05.1989"},
			Films:     []film.Film{{Title: "Film3"}, {Title: "Film4"}},
		},
	}

	// Calling the function
	result := film.ConvertMapToActorsListWithFilms(data)

	// Comparing the result with the expected value
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
