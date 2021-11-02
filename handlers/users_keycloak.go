package handlers

import (
	"context"
	"net/http"
	"token_based_auth/logger"
	"token_based_auth/middleware"

	"github.com/Nerzal/gocloak/v8"
	"github.com/gin-gonic/gin"
)

var AdminJWT *gocloak.JWT

func LoginUser(ctx *gin.Context) {
	AdminJWT, _ = middleware.AdminLogin()
	userInput := make(map[string]string)
	err := ctx.BindJSON(&userInput)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Inccorect user input",
		})
		return
	}

	gocloakClient := gocloak.NewClient("http://localhost:8080")
	jwt, err := gocloakClient.Login(context.Background(),
		middleware.KeyCloakClient,
		middleware.KeyCloakClientSecret,
		middleware.Realm,
		userInput["username"],
		userInput["password"],
	)

	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Inccorect username or password",
		})
		return
	}
	accesstoken := jwt.AccessToken
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Successfully logged in!!!!",
		"access_token": accesstoken,
	})
	logger.Write("success", "User logged in")

}

func CreateNewUser(ctx *gin.Context) {
	AdminJWT, _ = middleware.AdminLogin()
	newUser := make(map[string]interface{})
	err := ctx.BindJSON(&newUser)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect Create user format",
		})
		return
	}
	firstname := (newUser["firstname"].(string))
	lastname := (newUser["lastname"].(string))
	email := (newUser["email"].(string))
	username := (newUser["username"].(string))

	var goCloakUser gocloak.User
	goCloakUser.FirstName = &firstname
	goCloakUser.LastName = &lastname
	goCloakUser.Email = &email
	goCloakUser.Username = &username

	gocloakClient := gocloak.NewClient("http://localhost:8080")
	userId, err := gocloakClient.CreateUser(context.TODO(), AdminJWT.AccessToken, middleware.Realm, goCloakUser)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	latestUserDetails, err := gocloakClient.GetUserByID(context.TODO(), AdminJWT.AccessToken, middleware.Realm, userId)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"userdetails": latestUserDetails,
	})
	logger.Write("success", "New Keycloak user created")
}

func UpdateUser(ctx *gin.Context) {
	AdminJWT, _ = middleware.AdminLogin()
	updateData := make(map[string]interface{})
	err := ctx.BindJSON(&updateData)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect update user format",
		})
		return
	}
	firstname := (updateData["firstname"].(string))
	lastname := (updateData["lastname"].(string))
	email := (updateData["email"].(string))
	userID := (updateData["id"].(string))

	if userID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var goCloakUser gocloak.User
	goCloakUser.FirstName = &firstname
	goCloakUser.LastName = &lastname
	goCloakUser.Email = &email
	goCloakUser.ID = &userID

	gocloakClient := gocloak.NewClient("http://localhost:8080")
	err = gocloakClient.UpdateUser(context.TODO(), AdminJWT.AccessToken, middleware.Realm, goCloakUser)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	latestUserDetails, err := gocloakClient.GetUserByID(context.TODO(), AdminJWT.AccessToken, middleware.Realm, userID)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"userdetails": latestUserDetails,
	})
	logger.Write("success", "User data updated")
}

func DeleteUserByID(ctx *gin.Context) {
	AdminJWT, _ = middleware.AdminLogin()
	deleteData := make(map[string]interface{})
	err := ctx.BindJSON(&deleteData)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect delete user format",
		})
		return
	}

	userID := (deleteData["id"].(string))

	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect UserID",
		})
		return
	}

	gocloakClient := gocloak.NewClient("http://localhost:8080")
	err = gocloakClient.DeleteUser(context.TODO(), AdminJWT.AccessToken, middleware.Realm, userID)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, "User deleted Successfully! ")
	logger.Write("success", "User data deleted")
}

func GetAllUsers(ctx *gin.Context) {
	AdminJWT, _ = middleware.AdminLogin()
	gocloakClient := gocloak.NewClient("http://localhost:8080")
	var params gocloak.GetUsersParams
	users, err := gocloakClient.GetUsers(context.TODO(), AdminJWT.AccessToken, middleware.Realm, params)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
	logger.Write("success", "All users data fetched")
}
