package film

import (
	"encoding/json"
	"filmoteka/internal/actor"
	"filmoteka/pkg"
	"fmt"
	"log"
	"net/http"
)

type Storage interface {
	Add(film *Film) error
	GetFilmId(film *Film) (int64, error)
	Update(filmId int64, newFilm *Film) error
	Delete(filmId int64) error
	GetAllFilms(sortCol string) ([]Film, error)
	FindFilms(toFind string) ([]Film, error)
	ActorsListWithFilms() (map[actor.Actor][]Film, error)
}

type FilmHandler struct {
	FilmRepo Storage
}

// @Summary Добавляет фильм
// @Description Добавляет новый фильм в базу данных на основе переданных данных.
// @Accept json
// @Produce json
// @Param film body Film true "Данные фильма"
// @Success 201 {string} string "film added"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /user/film/add [post]
func (h *FilmHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	var film Film

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&film)
	if err != nil {
		log.Println("error decoding request JSON:", err)
		http.Error(w, "can't decode request JSON", http.StatusBadRequest)
		return
	}

	defer pkg.CloseBody(r)

	err = pkg.DateValidation(film.ReleaseDate)
	if err != nil {
		log.Println("wrong film release date format: ", err)
		http.Error(w, fmt.Sprintf("wrong film release date format"), http.StatusBadRequest)
		return
	}

	for _, actor := range film.Actors {
		err = pkg.DateValidation(actor.BirthDate)
		if err != nil {
			log.Println("wrong actor birth date format: %w", err)
			http.Error(w, fmt.Sprintf("wrong actor birth date format"), http.StatusBadRequest)
			return
		}
	}

	err = h.FilmRepo.Add(&film)
	if err != nil {
		log.Println("error adding film:", err)
		http.Error(w, "can't add film: film already exists", http.StatusInternalServerError)
		return
	}

	log.Println("film added:", film)
	w.Write([]byte(fmt.Sprintf("film added: %v", film)))
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновляет информацию о фильме
// @Description Обновляет информацию о фильме в базе данных на основе переданных данных.
// @Accept json
// @Produce json
// @Param filmInfo body []Film true "Старая и новая информация о фильме"
// @Success 200 {string} string "film updated"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/film/update [put]
func (h *FilmHandler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	var filmInfo []Film

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&filmInfo)
	if err != nil {
		log.Println("error decoding request JSON:", err)
		http.Error(w, "can't decode request JSON", http.StatusBadRequest)
		return
	}

	oldFilm := filmInfo[0]
	newFilm := filmInfo[1]

	defer pkg.CloseBody(r)

	err = pkg.DateValidation(newFilm.ReleaseDate)
	if err != nil {
		log.Println("wrong film release date format: ", err)
		http.Error(w, fmt.Sprintf("wrong film release date format"), http.StatusBadRequest)
		return
	}

	for _, actor := range newFilm.Actors {
		err = pkg.DateValidation(actor.BirthDate)
		if err != nil {
			log.Println("wrong actor birth date format: %w", err)
			http.Error(w, fmt.Sprintf("wrong actor birth date format"), http.StatusBadRequest)
			return
		}
	}

	oldFilmId, err := h.FilmRepo.GetFilmId(&oldFilm)
	if err != nil {
		log.Println("error getting film id:", err)
		http.Error(w, "can't find film", http.StatusInternalServerError)
		return
	}

	err = h.FilmRepo.Update(oldFilmId, &newFilm)
	if err != nil {
		log.Println("error updating film:", err)
		http.Error(w, "can't update film", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("film updated"))
	w.WriteHeader(http.StatusOK)
}

// @Summary Удаляет фильм
// @Description Удаляет фильм из базы данных на основе переданных данных.
// @Accept json
// @Produce json
// @Param film body Film true "Данные фильма"
// @Success 200 {string} string "film deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/film/delete [delete]
func (h *FilmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	var film Film

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&film)
	if err != nil {
		log.Println("error decoding request JSON:", err)
		http.Error(w, "can't decode request JSON", http.StatusBadRequest)
		return
	}

	defer pkg.CloseBody(r)

	err = film.Validate(w)
	if err != nil {
		return
	}

	filmId, err := h.FilmRepo.GetFilmId(&film)
	if err != nil {
		log.Println("error getting film id:", err)
		http.Error(w, "can't find film", http.StatusInternalServerError)
		return
	}

	err = h.FilmRepo.Delete(filmId)
	if err != nil {
		log.Println("error deleting film:", err)
		http.Error(w, "can't delete film", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("film deleted"))
	w.WriteHeader(http.StatusOK)
}

// @Summary Получает список всех фильмов
// @Description Возвращает список всех фильмов из базы данных, с возможностью сортировки по указанному столбцу (по умолчанию сортировка по названию).
// @Produce json
// @Param sort query string false "Столбец для сортировки (title, release_date, rating)"
// @Success 200 {array} Film "Список фильмов"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /user/films [get]
func (h *FilmHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	sortCol := r.URL.Query().Get("sort")
	fmt.Println(sortCol)
	if sortCol == "" {
		sortCol = "title"
	}
	if sortCol != "title" && sortCol != "release_date" && sortCol != "rating" {
		log.Println("error getting all films: wrong sort column")
		http.Error(w, "wrong sort column: it must be empty, title, release_date or rating", http.StatusBadRequest)
		return
	}

	films, err := h.FilmRepo.GetAllFilms(sortCol)
	if err != nil {
		log.Println("error getting all films:", err)
		http.Error(w, "can't get all films", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(films)
	if err != nil {
		log.Println("error marshalling films:", err)
		http.Error(w, "can't marshal films", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

// @Summary Находит фильмы по строке поиска
// @Description Поиск фильмов в базе данных по указанной строке поиска.
// @Produce json
// @Param find query string true "Строка поиска"
// @Success 200 {array} Film "Найденные фильмы"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /user/films/find [get]
func (h *FilmHandler) FindFilms(w http.ResponseWriter, r *http.Request) {
	toFind := r.URL.Query().Get("find")
	if toFind == "" {
		log.Println("error finding films: empty find string")
		http.Error(w, "empty find string", http.StatusBadRequest)
		return
	}

	films, err := h.FilmRepo.FindFilms(toFind)
	if err != nil {
		log.Println("error finding films:", err)
		http.Error(w, "can't find films", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(films)
	if err != nil {
		log.Println("error marshalling films:", err)
		http.Error(w, "can't marshal films", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

// @Summary Получает список актеров с их фильмами
// @Description Возвращает список всех актеров вместе с их фильмами из базы данных.
// @Produce json
// @Success 200 {array} ActorListWithFilms "Список актеров с фильмами"
// @Failure 500 {string} string "Internal server error"
// @Router /user/actors [get]
func (h *FilmHandler) ActorsListWithFilms(w http.ResponseWriter, r *http.Request) {
	actorsList, err := h.FilmRepo.ActorsListWithFilms()
	if err != nil {
		log.Println("error getting actors list with films:", err)
		http.Error(w, "can't get actors list with films", http.StatusInternalServerError)
		return
	}
	actorsListWithFilms := ConvertMapToActorsListWithFilms(actorsList)
	resp, err := json.Marshal(actorsListWithFilms)
	if err != nil {
		log.Println("error marshalling actors list with films:", err)
		http.Error(w, "can't marshal actors list with films", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
