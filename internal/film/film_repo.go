package film

import (
	"database/sql"
	"filmoteka/internal/actor"
	"fmt"
)

type FilmRepository struct {
	actorRepo *actor.ActorRepository
	db        *sql.DB
}

func NewFilmRepository(actorRepo *actor.ActorRepository, db *sql.DB) *FilmRepository {
	return &FilmRepository{
		actorRepo: actorRepo,
		db:        db,
	}
}

func (repo *FilmRepository) Add(film *Film) error {
	op := "film_repo.Add"

	row := repo.db.QueryRow("INSERT INTO film(title, description, release_date, rating) VALUES($1, $2, $3, $4) RETURNING id",
		film.Title, film.Description, film.ReleaseDate, film.Rating)

	var filmId int64
	err := row.Scan(&filmId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	for _, actor := range film.Actors {
		actorId, err := repo.actorRepo.GetActorId(&actor)
		if err != nil {
			repo.actorRepo.Add(&actor)
			actorId, err = repo.actorRepo.GetActorId(&actor)
		}
		_, err = repo.db.Exec("INSERT INTO film_actor(film_id, actor_id) VALUES($1, $2)", filmId, actorId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (repo *FilmRepository) GetFilmId(film *Film) (int64, error) {
	op := "film_repo.GetByID"
	row := repo.db.QueryRow("SELECT id FROM film where title = $1 and release_date = $2",
		film.Title, film.ReleaseDate)
	var filmId int64
	err := row.Scan(&filmId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return filmId, nil
}

func (repo *FilmRepository) Update(filmId int64, newFilm *Film) error {
	op := "film_repo.UpdateFilm"
	_, err := repo.db.Query("UPDATE film SET title = $1, description = $2, release_date = $3, rating = $4 WHERE id = $5",
		newFilm.Title, newFilm.Description, newFilm.ReleaseDate, newFilm.Rating, filmId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (repo *FilmRepository) Delete(filmId int64) error {
	op := "film_repo.DeleteFilm"
	tx, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec("DELETE FROM film_actor WHERE film_id = $1", filmId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = tx.Exec("DELETE FROM film WHERE id = $1", filmId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (repo *FilmRepository) GetAllFilms(sortCol string) ([]Film, error) {
	op := "film_repo.GetAllFilms"

	var validCols = map[string]bool{
		"title":        true,
		"release_date": true,
		"rating":       true,
	}
	if !validCols[sortCol] {
		return nil, fmt.Errorf("%s: invalid column name: %s", op, sortCol)
	}

	stmt, err := repo.db.Prepare(fmt.Sprintf("SELECT title, description, release_date, rating FROM film ORDER BY %s", sortCol))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.Title, &film.Description, &film.ReleaseDate, &film.Rating)
		film.Actors = nil
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		films = append(films, film)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return films, nil
}

func (repo *FilmRepository) FindFilms(toFind string) ([]Film, error) {
	op := "film_repo.FindFilm"

	var films []Film
	var err error
	if films, err = repo.findFilmsByTitle(toFind); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(films) == 0 {
		if films, err = repo.findFilmsByActor(toFind); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if len(films) == 0 {
		return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("no films found"))
	}

	return films, nil
}

func (repo FilmRepository) ActorsListWithFilms() (map[actor.Actor][]Film, error) {
	op := "film_repo.ActorListWithFilms"

	rows, err := repo.db.Query(`SELECT distinct name, gender, birth_date from actor`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	actorsWithFilms := make(map[actor.Actor][]Film)

	for rows.Next() {
		var actor actor.Actor
		err := rows.Scan(&actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actorFilms, err := repo.findFilmsByActor(actor.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actorsWithFilms[actor] = actorFilms
	}
	return actorsWithFilms, nil
}

func (repo *FilmRepository) findFilmsByTitle(titleFragment string) ([]Film, error) {
	op := "film_repo.FindFilmsByTitleFragment"

	query := fmt.Sprintf(`
        SELECT f.title, f.description, f.release_date, f.rating
        FROM film f
        WHERE f.title LIKE '%%%s%%'`, titleFragment)

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.Title, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		films = append(films, film)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return films, nil
}

func (repo *FilmRepository) findFilmsByActor(actorName string) ([]Film, error) {
	op := "film_repo.FindFilmsByActor"

	query := fmt.Sprintf(`
    SELECT f.title, f.description, f.release_date, f.rating
    FROM film f
    WHERE f.id IN (
        SELECT film_id 
        FROM film_actor 
        WHERE actor_id IN (
            SELECT id 
            FROM actor 
            WHERE name LIKE '%%%s%%'
        )
    )`, actorName)

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		err := rows.Scan(&film.Title, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		films = append(films, film)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return films, nil
}
