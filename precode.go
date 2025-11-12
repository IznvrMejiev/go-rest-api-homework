package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{}

func getTasks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {

		http.Error(w, "Ошибка при кодировании задач", http.StatusInternalServerError)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {

		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if newTask.ID == "" {
		http.Error(w, "ID обязателен", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	taskID := chi.URLParam(r, "id")

	task, exists := tasks[taskID]
	if !exists {

		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	taskID := chi.URLParam(r, "id")

	_, exists := tasks[taskID]
	if !exists {

		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Задача удалена"})
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)

	r.Post("/tasks", createTask)

	r.Get("/tasks/{id}", getTask)

	r.Delete("/tasks/{id}", deleteTask)

	fmt.Println("✅ Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
