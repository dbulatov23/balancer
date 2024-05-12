package main

import (
	"balancer/pkg/api"
	"balancer/pkg/users"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	env := os.Getenv("ADDRESS")
	fmt.Println("starting service", env)

	s := grpc.NewServer()
	srv := &users.UserInfo{}
	api.RegisterUserServer(s, srv)

	l, err := net.Listen("tcp", env)
	if err != nil {
		log.Fatal(err)

	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
