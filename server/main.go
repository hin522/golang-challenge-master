package main

import (
	"log"
	"main/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Server is starting...")

	r := gin.Default()

	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)
	r.POST("/users/messages", handlers.CreateMessage)
	r.PATCH("/users/remove-nickname", handlers.RemoveUserNickname)

	r.Run()
}
