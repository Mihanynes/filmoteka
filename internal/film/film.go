package film

import (
	"filmoteka/internal/actor"
	"filmoteka/pkg"
	"fmt"
	"log"
	"net/http"
)

type Film struct {
	Title       string        `json:"title" notempty:"true"`
	Description string        `json:"description" notempty:"true"`
	ReleaseDate string        `json:"release_date" notempty:"true" validate:"date"`
	Rating      int           `json:"rating" notempty:"true" validate:"min=1,max=10"`
	Actors      []actor.Actor `json:"actors,omitempty"`
}

type ActorListWithFilms struct {
	ActorInfo actor.Actor `json:"actor"`
	Films     []Film      `json:"films"`
}

func (film *Film) Validate(w http.ResponseWriter) error {
	err := pkg.DateValidation(film.ReleaseDate)
	if err != nil {
		log.Println("wrong film release date format: ", err)
		http.Error(w, fmt.Sprintf("wrong film release date format"), http.StatusBadRequest)
		return err
	}

	for _, actor := range film.Actors {
		err = pkg.DateValidation(actor.BirthDate)
		if err != nil {
			log.Println("wrong actor birth date format: %w", err)
			http.Error(w, fmt.Sprintf("wrong actor birth date format"), http.StatusBadRequest)
			return err
		}
	}
	return nil
}

func ConvertMapToActorsListWithFilms(data map[actor.Actor][]Film) []ActorListWithFilms {
	var result []ActorListWithFilms

	for actor, films := range data {
		actorWithFilms := ActorListWithFilms{
			ActorInfo: actor,
			Films:     films,
		}
		result = append(result, actorWithFilms)
	}

	return result
}
