package actor

import (
	"encoding/json"
	"filmoteka/pkg"
	"fmt"
	"log"
	"net/http"
)

type Storage interface {
	Add(*Actor) error
	GetActorId(*Actor) (int64, error)
	Update(int64, *Actor) error
	Delete(int64) error
}

type ActorHandler struct {
	ActorRepo Storage
}

// @Summary Добавляет актера
// @Description Добавляет нового актера в базу данных на основе переданных данных.
// @Accept json
// @Produce json
// @Param actor body Actor true "Данные актера"
// @Success 201 {string} string "actor added: {actor}"
// @Failure 400 {string} string "Bad request"
// @Failure 409 {string} string "Actor already exists"
// @Failure 500 {string} string "Internal server error"
// @Router user/actor/add [post]
func (h *ActorHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	var actor Actor

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&actor)
	if err != nil {
		log.Println("error decoding request JSON:", err)
		http.Error(w, "can't decode request JSON", http.StatusBadRequest)
		return
	}

	defer pkg.CloseBody(r)

	err = pkg.DateValidation(actor.BirthDate)
	if err != nil {
		log.Println("wrong actor birth date format: ", err)
		http.Error(w, fmt.Sprintf("wrong actor birth date format"), http.StatusBadRequest)
		return
	}

	_, err = h.ActorRepo.GetActorId(&actor)
	if err == nil {
		log.Println("error adding actor: actor already exists", err)
		http.Error(w, "actor already exists", http.StatusConflict)
		return
	}

	err = h.ActorRepo.Add(&actor)
	if err != nil {
		log.Println("error adding actor:", err)
		http.Error(w, fmt.Sprintf("can't add actor %w", err), http.StatusInternalServerError)
		return
	}

	log.Println("actor added:", actor)
	w.Write([]byte(fmt.Sprintf("actor added: %v", actor)))
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновляет информацию об актере
// @Description Обновляет информацию об актере в базе данных на основе переданных данных.
// @Accept json
// @Produce json
// @Param actorInfo body []Actor true "Старая и новая информация об актере"
// @Success 200 {string} string "actor updated: {newActor}"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Actor not found"
// @Failure 500 {string} string "Internal server error"
// @Router admin/actor/update [put]
func (h *ActorHandler) UpdateActor(w http.ResponseWriter, r *http.Request) {

	actorInfo := make([]Actor, 2)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&actorInfo)

	if err != nil {
		log.Println("error decoding request JSON:", err)
		http.Error(w, "can't decode request JSON", http.StatusBadRequest)
		return
	}

	oldActor := actorInfo[0]
	newActor := actorInfo[1]

	defer pkg.CloseBody(r)

	err = pkg.DateValidation(newActor.BirthDate)
	if err != nil {
		log.Println("wrong actor birth date format: ", err)
		http.Error(w, fmt.Sprintf("wrong actor birth date format"), http.StatusBadRequest)
		return
	}

	oldActorID, err := h.ActorRepo.GetActorId(&oldActor)
	if err != nil {
		log.Println("error updating actor: actor not found", err)
		http.Error(w, "actor not found", http.StatusNotFound)
		return
	}

	err = h.ActorRepo.Update(oldActorID, &newActor)
	if err != nil {
		log.Println("error updating actor:", err)
		http.Error(w, fmt.Sprintf("can't update actor %w", err), http.StatusInternalServerError)
		return
	}

	log.Println("actor updated:", newActor)
	w.Write([]byte(fmt.Sprintf("actor updated: %v", newActor)))
	w.WriteHeader(http.StatusOK)
}

// @Summary Удаляет актера
// @Description Удаляет актера из базы данных по переданным данным актера.
// @Accept json
// @Produce json
// @Param actor body Actor true "Данные актера"
// @Success 200 {string} string "actor deleted: {actor}"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Actor not found"
// @Failure 500 {string} string "Internal server error"
// @Router admin/actor/delete [delete]
func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	var actor Actor

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&actor)
	if err != nil {
		log.Println("error decoding request JSON:", err)
		http.Error(w, "can't decode request JSON", http.StatusBadRequest)
		return
	}

	defer pkg.CloseBody(r)

	err = pkg.DateValidation(actor.BirthDate)
	if err != nil {
		log.Println("wrong actor birth date format: ", err)
		http.Error(w, fmt.Sprintf("wrong actor birth date format"), http.StatusBadRequest)
		return
	}

	actorID, err := h.ActorRepo.GetActorId(&actor)
	if err != nil {
		log.Println("error deleting actor: actor not found", err)
		http.Error(w, "actor not found", http.StatusNotFound)
		return
	}

	err = h.ActorRepo.Delete(actorID)
	if err != nil {
		log.Println("error deleting actor:", err)
		http.Error(w, fmt.Sprintf("can't delete actor %w", err), http.StatusInternalServerError)
		return
	}

	log.Println("actor deleted:", actor)
	w.Write([]byte(fmt.Sprintf("actor deleted: %v", actor)))
	w.WriteHeader(http.StatusOK)
}
