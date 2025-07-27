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
	MessageCount int `json:"messsageCount"`
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
			MessageCount: row.MessageCount,
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

type CreateMessageRequest struct {
	Username string `json:"username" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

type CreateMessageResponse struct {
	ID int `json:"id"`
}

func CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	id, err := queries.InsertMessage(queries.InsertMessageParams{
		Username: req.Username,
		Message:  req.Message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create message: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateMessageResponse{ID: id})
}

type RemoveNicknameRequest struct {
	Username string `json:"username" binding:"required"`
}

func RemoveUserNickname(c *gin.Context) {
	var req RemoveNicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	err := queries.RemoveUserNickname(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove nickname: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Nickname removed successfully",
	})
}
