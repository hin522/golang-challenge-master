package handlers

import (
	"main/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type UserResponse struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	UserType string  `json:"userType"`
	Nickname *string `json:"nickname,omitempty"`
}

func GetUsers(c *gin.Context) {

	userRows, err := queries.GetUsers()

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to retrieve users: " + err.Error(),
		})
		return
	}

	var users []UserResponse = []UserResponse{}
	for _, row := range userRows {

		var nickname *string = nil 
		if row.Nickname.Valid {
			nickname = &row.Nickname.String
		}

		users = append(users, UserResponse{
			ID:       row.ID,
			Username: row.Username,
			Email:    row.Email,
			UserType: row.UserType,
			Nickname: nickname,
		})
	}

	c.JSON(http.StatusOK, GetUsersResponse{
		Users: users,
	})
}

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required"`
	UserType string  `json:"userType" binding:"required"`
	Nickname *string `json:"nickname,omitempty"`
}

type CreateUserResponse struct {
	ID int `json:"id"`
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	// Call the CreateUser query
	params := queries.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		UserType: req.UserType,
		Nickname: "",
	}

	if req.Nickname != nil {
		params.Nickname = *req.Nickname
	}

	id, err := queries.CreateUser(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateUserResponse{ID: id})
}