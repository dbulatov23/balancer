package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	City      string `json:"city"`
}

var users = []User{
	{
		ID:        1,
		Key:       "6342ff6e-b2de-4059-a19c-389bb1f79e3a",
		FirstName: "Даниял",
		LastName:  "Булатов",
		City:      "Москва",
	},
	{
		ID:        2,
		Key:       "bfb0e275-0a93-4f54-8d5b-b5a569ed7647",
		FirstName: "Иван",
		LastName:  "Иванов",
		City:      "Казань",
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func getUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.Host)
	w.Write(resp)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user_id, _ := strconv.Atoi(id)
	var cnt, ind int
	for i, v := range users {
		if v.ID == user_id {
			cnt++
			ind = i
		} else {
			continue
		}
	}
	if cnt == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(users[ind])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func createUsers(w http.ResponseWriter, r *http.Request) {
	var user User
	var buf bytes.Buffer
	var cnt_key int

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	for _, v := range users {
		fmt.Println(user.Key, v.Key)
		if user.Key == v.Key {
			cnt_key = 1
			break
		}
	}
	if cnt_key != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("already exist uuid"))
	} else {
		users = append(users, user)
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	r := chi.NewRouter()
	// здесь регистрируйте ваши обработчики
	r.Get("/users", getUsers)
	r.Post("/users", createUsers)
	r.Get("/users/{id}", getUser)

	env := os.Getenv("ADDRESS")
	fmt.Println("starting service", env)
	if err := http.ListenAndServe(env, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
