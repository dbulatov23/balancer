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
	"time"

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

func proxy(w http.ResponseWriter, r *http.Request) {
	var req *api.GetUsersRequest
	for i := 0; i < 3; i++ {
		fmt.Println(r.URL.Path)
		host := getHost()
		conn, err := grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		c := api.NewUserClient(conn)
		if r.URL.Path == "/users" && r.Method == "GET" {
			res, err := c.GetUsers(context.Background(), req)
			if err != nil {
				log.Fatal(err)
				time.Sleep(3 * time.Second)
				continue
			} else {
				message, err := json.Marshal(res)
				if err != nil {
					log.Fatal(err)
				}
				w.Write(message)
				w.WriteHeader(http.StatusOK)
				break
			}
		} else if r.URL.Path == "/user" && r.Method == "GET" {
			params := r.URL.Query()
			param := params.Get("id")
			num, err := strconv.Atoi(param)
			if err != nil {
				fmt.Println("Ошибка преобразования:", err)
				return
			}
			res, err := c.GetUser(context.Background(), &api.GetUserRequest{Id: int32(num)})
			if err != nil {
				log.Fatal(err)
				time.Sleep(3 * time.Second)
				continue
			} else {
				message, err := json.Marshal(res)
				if err != nil {
					log.Fatal(err)
				}
				w.Write(message)
				w.WriteHeader(http.StatusOK)
				break
			}
		} else if r.URL.Path == "/user" && r.Method == "POST" {
			var user User
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal(body, &user)
			res, err := c.CreateUser(context.Background(), &api.UserInfo{
				Id:        int32(user.ID),
				Key:       user.Key,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				City:      user.City,
			})
			if err != nil {
				log.Fatal(err)
				time.Sleep(3 * time.Second)
				continue
			} else {
				message, err := json.Marshal(res)
				if err != nil {
					log.Fatal(err)
				}
				w.Write(message)
				w.WriteHeader(http.StatusOK)
				break
			}
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
	r.Get("/user", proxy)
	if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
