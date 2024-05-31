package users

import (
	"balancer/pkg/api"
	b "balancer/pkg/api"
	"context"
	"errors"
)

type UserInfo struct {
	b.UnimplementedUserServer
}

type User struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	City      string `json:"city"`
}

var ErrUserNotFound = errors.New("User not found")

var users = []User{
	{
		ID:        1,
		Key:       "6342ff6e-b2de-4059-a19c-389bb1f79e3a",
		FirstName: "Анвар",
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

func (*UserInfo) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.UserInfo, error) {
	ind := req.GetId()
	var cnt int
	for _, v := range users {
		if v.ID == int(ind) {
			cnt++
		} else {
			continue
		}
	}
	if cnt == 0 {
		return &api.UserInfo{}, nil
	}
	x := users[req.GetId()-1]
	return &api.UserInfo{
		Id:        int32(x.ID),
		Key:       x.Key,
		FirstName: x.FirstName,
		LastName:  x.LastName,
		City:      x.City,
	}, nil
}

func (*UserInfo) GetUsers(ctx context.Context, req *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	response := make([]*api.UserInfo, len(users))
	for i, v := range users {
		response[i] = &api.UserInfo{
			Id:        int32(v.ID),
			Key:       v.Key,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			City:      v.City,
		}
	}
	return &api.GetUsersResponse{
		UsersInfo: response,
	}, nil
}

func (*UserInfo) CreateUser(ctx context.Context, req *api.UserInfo) (*api.CreateUserResponse, error) {
	var str *api.CreateUserResponse
	var cnt_key int
	for _, v := range users {
		if req.GetKey() == v.Key {
			cnt_key = 1
			break
		}
	}
	if cnt_key != 0 {
		return str, nil
	}
	users = append(users, User{
		ID:        int(req.GetId()),
		Key:       req.GetKey(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		City:      req.City,
	})
	return str, nil
}
