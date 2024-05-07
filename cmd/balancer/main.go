package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

var Id int

func roundRobin() string {
	Id = (Id + 1) % 3
	steps := viper.GetStringSlice("steps")
	return steps[Id]
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 3; i++ {
		host := roundRobin()
		response, err := http.Get("http://" + host + "/users")
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		} else {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Write(body)
			w.WriteHeader(http.StatusOK)
			break
		}
		w.WriteHeader(http.StatusBadGateway)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 3; i++ {
		host := roundRobin()
		response, err := http.Get("http://" + host + "/user")
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		} else {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Write(body)
			w.WriteHeader(http.StatusOK)
			break
		}
		w.WriteHeader(http.StatusBadGateway)
	}
}

func createUsers(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 3; i++ {
		host := roundRobin()
		fmt.Println(host)
		response, err := http.Post("http://"+host+"/users", "application/json", r.Body)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.StatusCode)
			body, _ := io.ReadAll(response.Body)
			w.Write(body)
			break
		}
		w.WriteHeader(http.StatusBadGateway)
	}
}

// Ниже напишите обработчики для каждого эндпоинта

func main() {
	r := chi.NewRouter()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	r.Get("/users", getUsers)
	r.Post("/users", createUsers)
	r.Get("/users/{id}", getUser)
	// здесь регистрируйте ваши обработчики
	if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
