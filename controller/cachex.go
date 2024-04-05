package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"modulgo/model"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetUsers() []model.User {
	client := getClient()
	val, err := client.Get(ctx, "users").Result()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Check Redis")
	var users []model.User
	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		log.Println(err)
	}
	return users
}

func SetUsers(users []model.User) {
	converted, err := json.Marshal(users)
	if err != nil {
		log.Println("JSON Marshal Error: ", err)
	}
	client := getClient()
	defer client.Close()
	if err := client.Set(ctx, "users", converted, 0).Err(); err != nil {
		log.Println("Redis Set Error: ", err)
	} else {
		log.Println("Redis Set Success")
	}
}
