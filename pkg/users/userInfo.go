package users

import (
	"balancer/pkg/api"
	b "balancer/pkg/api"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

type UserCache struct {
	Id   string
	City interface{}
	TTL  time.Duration
}

type UserInfo struct {
	b.UnimplementedUserServer
	DB    *sql.DB
	Redis *redis.Client
}
type User struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	City      string `json:"city"`
}

var ErrUserNotFound = errors.New("User not found")

func NewUserInfo(db *sql.DB, redis *redis.Client) *UserInfo {
	return &UserInfo{
		DB:    db,
		Redis: redis,
	}
}

func (u *UserInfo) GetUsers(ctx context.Context, req *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	rows, err := u.DB.Query("SELECT id, key, first_name, last_name, city FROM public.user")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Key, &user.FirstName, &user.LastName, &user.City)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

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

func (u *UserInfo) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.UserInfo, error) {

	ind := req.GetId()
	data, err := u.Redis.Get(string(ind)).Result()
	if err != nil {
		fmt.Println(err)
		fmt.Println("ошибка")
	}
	if data != "" {
		return &api.UserInfo{
			Id:        ind,
			Key:       "из кэша",
			FirstName: "из кэша",
			LastName:  "из кэша",
			City:      data,
		}, nil
	}
	if data == "" {
		fmt.Println("not data in redis")
	}

	row := u.DB.QueryRow("SELECT id, key, first_name, last_name, city FROM public.user where id = $1", ind)
	var user User

	err = row.Scan(&user.ID, &user.Key, &user.FirstName, &user.LastName, &user.City)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(&user.City)
	if err != nil {
		return nil, err
	}
	err = u.Redis.Set(string(ind), b, 10*time.Second).Err()
	if err != nil {
		fmt.Println("Ошибка при добавлении данных в кэш:", err)
		return nil, err
	}
	fmt.Println("Записал в редис")

	return &api.UserInfo{
		Id:        int32(user.ID),
		Key:       user.Key,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		City:      user.City,
	}, nil
}
