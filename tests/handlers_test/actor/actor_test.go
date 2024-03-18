package actorTest

import (
	"bytes"
	"encoding/json"
	"filmoteka/internal/actor"
	"fmt"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestActorHandler_AddActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	handler := &actor.ActorHandler{
		ActorRepo: mockStorage,
	}

	// Успешное добавление нового актера
	testActor := &actor.Actor{
		Name:      "John",
		BirthDate: "01.01.1990",
	}
	mockStorage.EXPECT().GetActorId(testActor).Return(int64(0), fmt.Errorf("actor not found"))
	mockStorage.EXPECT().Add(testActor).Return(nil)

	reqBody, err := json.Marshal(testActor)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/actor", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.AddActor(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedResponse := "actor added: {John  01.01.1990}"
	if body := strings.TrimSpace(w.Body.String()); body != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, body)
	}
}

func TestActorHandler_UpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	handler := &actor.ActorHandler{
		ActorRepo: mockStorage,
	}

	// Создаем тестовых актеров для обновления
	oldActor := &actor.Actor{Name: "John", BirthDate: "01.01.1990"}
	newActor := &actor.Actor{Name: "John", BirthDate: "02.02.1990"}

	// Устанавливаем ожидаемое поведение мока GetActorId
	mockStorage.EXPECT().GetActorId(oldActor).Return(int64(1), nil)

	// Устанавливаем ожидаемое поведение мока Update
	mockStorage.EXPECT().Update(int64(1), newActor).Return(nil)

	// Создаем JSON-данные для обновления актера
	actorInfo := []actor.Actor{*oldActor, *newActor}
	reqBody, err := json.Marshal(actorInfo)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/actors", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.UpdateActor(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedResponse := "actor updated: {John  02.02.1990}"
	if body := strings.TrimSpace(w.Body.String()); body != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, body)
	}
}

func TestActorHandler_DeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := NewMockStorage(ctrl)
	handler := &actor.ActorHandler{
		ActorRepo: mockStorage,
	}

	// Создаем тестового актера для удаления
	testActor := &actor.Actor{Name: "John", BirthDate: "01.01.1990"}

	mockStorage.EXPECT().GetActorId(testActor).Return(int64(1), nil)

	mockStorage.EXPECT().Delete(int64(1)).Return(nil)

	reqBody, err := json.Marshal(testActor)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest("POST", "/actors", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.DeleteActor(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedResponse := "actor deleted: {John  01.01.1990}"
	if body := strings.TrimSpace(w.Body.String()); body != expectedResponse {
		t.Errorf("expected response body %q, got %q", expectedResponse, body)
	}
}
