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

func getHost() string {
	Id = (Id + 1) % 3
	upstreams := viper.GetStringSlice("upstreams")
	return upstreams[Id]
}

func proxy(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	method := r.Method
	for i := 0; i < 3; i++ {
		host := getHost()
		response, err := http.NewRequest(method, "http://"+host+uri, nil)
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
	}
}

func main() {
	r := chi.NewRouter()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	r.Get("/users", proxy)
	r.Post("/user", proxy)
	r.Get("/users/{id}", proxy)
	// здесь регистрируйте ваши обработчики
	if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
