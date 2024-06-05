package main

import (
	"balancer/pkg/api"
	"balancer/pkg/users"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"
)

func main() {
	env := os.Getenv("ADDRESS")
	fmt.Println("starting service", env)
	connStr := "user=postgres password=1234 dbname=postgres sslmode=disable host=postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	client := redis.NewClient(&redis.Options{
		Addr: "some-redis:6379",
	})

	s := grpc.NewServer()
	srv := users.NewUserInfo(db, client)
	api.RegisterUserServer(s, srv)

	l, err := net.Listen("tcp", env)
	if err != nil {
		log.Fatal(err)

	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
