package main

import (
	"balancer/pkg/api"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var Id int

type User struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	City      string `json:"city"`
}

func getHost() string {
	Id = (Id + 1) % 3
	upstreams := viper.GetStringSlice("upstreams")
	return upstreams[Id]
}

func getUsers(c api.UserClient) (string, error) {
	var req *api.GetUsersRequest
	res, err := c.GetUsers(context.Background(), req)
	if err != nil {
		return "", err
	}
	message, err := json.Marshal(res)
	return string(message), nil
}

func getUser(id string, c api.UserClient) (string, error) {
	num, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	res, err := c.GetUser(context.Background(), &api.GetUserRequest{Id: int32(num)})
	if err != nil {
		return "", err

	}
	message, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	return string(message), nil
}

func createUser(body []byte, c api.UserClient) (string, error) {
	var user User
	json.Unmarshal(body, &user)
	res, err := c.CreateUser(context.Background(), &api.UserInfo{
		Id:        int32(user.ID),
		Key:       user.Key,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		City:      user.City,
	})
	if err != nil {
		return "", err
	}
	message, err := json.Marshal(res)
	return string(message), nil
}

func proxy(w http.ResponseWriter, r *http.Request) {
	host := getHost()
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := api.NewUserClient(conn)
	if r.URL.Path == "/users" && r.Method == "GET" {
		res, err := getUsers(c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write([]byte(res))
	} else if r.URL.Path == "/user" && r.Method == "GET" {
		id := r.URL.Query().Get("id")
		res, err := getUser(id, c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(res))
		w.WriteHeader(http.StatusOK)
	} else if r.URL.Path == "/user" && r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		res, err := createUser(body, c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		message, err := json.Marshal(res)
		w.Write([]byte(message))
		w.WriteHeader(http.StatusOK)
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
	r.Get("/user", proxy)
	if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
