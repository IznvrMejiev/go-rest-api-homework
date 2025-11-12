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

// getTasks — GET /tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "Ошибка при кодировании задач", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
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

	data, err := json.Marshal(newTask)
	if err != nil {
		http.Error(w, "Ошибка при кодировании задачи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	taskID := chi.URLParam(r, "id")
	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Ошибка при кодировании задачи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
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

	response := map[string]string{"message": "Задача удалена"}
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	fmt.Println("✅ Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err.Error())
		return
	}
}
