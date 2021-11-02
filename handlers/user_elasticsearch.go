package handlers

import (
	"encoding/json"
	"net/http"
	"token_based_auth/logger"
	"token_based_auth/middleware"
	"token_based_auth/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNewUserES(ctx *gin.Context) {
	//check authourization
	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "u.c")

	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to create new users")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to create users",
		})
		return
	}
	logger.Write("success", "User Authorized to create new users")

	//get user details
	newUser := make(map[string]interface{})
	err := ctx.BindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Inccorect insert data format",
		})
		return
	}

	newUser["id"] = uuid.Must(uuid.NewRandom()).String()
	//insert in the database
	byteData, err := json.Marshal(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occured while creating the user",
		})
		return
	}

	x, err := model.CreateNewUser(string(byteData))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occured while creating the user",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"data":    x,
		"user_id": newUser["id"],
	})
	logger.Write("success", "New user created")
}

func GetAllUsersES(ctx *gin.Context) {
	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "u.r")

	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to read user details")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to read users",
		})
		return
	}
	logger.Write("success", "User Authorized to read all users data")

	var result []model.User
	err := model.GetAllUsers(&result)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occured while fetching all users data",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, result)
	logger.Write("success", "All user data fetched")
}

func UpdateUserES(ctx *gin.Context) {
	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "u.u")

	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to update user details")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to udpate users",
		})
		return
	}
	logger.Write("success", "User Authorized to update users data")

	uID := ctx.DefaultQuery("userID", "")
	if uID == "" {
		logger.Write("error", "Invalid user id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	updateData := make(map[string]interface{})
	err := ctx.BindJSON(&updateData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Inccorect update data format",
		})
		return
	}

	err, message := model.UpdateUserByID(uID, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user_id": uID,
		"message": "User has been updated successfully",
	})
	logger.Write("success", "User data updated")
}

func DeleteUsersES(ctx *gin.Context) {
	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "u.d")

	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to delete users")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to delete users",
		})
		return
	}
	logger.Write("success", "User Authorized to delete users")

	uID := ctx.DefaultQuery("userID", "")
	if uID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
	}

	err, message := model.DeleteUserByID(uID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id": uID,
		"message": "User has been deleted successfully",
	})
	logger.Write("success", "User data deleted")
}
