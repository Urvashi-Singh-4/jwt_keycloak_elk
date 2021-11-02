package handlers

import (
	"context"
	"testing"
	"token_based_auth/middleware"

	"github.com/Nerzal/gocloak/v8"
)

func BenchmarkGetAllUsers(b *testing.B) {
	AdminJWT, _ = middleware.AdminLogin()
	gocloakClient := gocloak.NewClient("http://localhost:8080")

	for i := 0; i < b.N; i++ {
		var params gocloak.GetUsersParams
		gocloakClient.GetUsers(context.TODO(), AdminJWT.AccessToken, middleware.Realm, params)
	}
}

func BenchmarkCreateandDeleteUsers(b *testing.B) {

	AdminJWT, _ = middleware.AdminLogin()

	firstname := "Tony"
	lastname := "Stark"
	email := "TonyStark@marvel.com"
	username := "iAmIronMan"

	var goCloakUser gocloak.User
	goCloakUser.FirstName = &firstname
	goCloakUser.LastName = &lastname
	goCloakUser.Email = &email
	goCloakUser.Username = &username

	gocloakClient := gocloak.NewClient("http://localhost:8080")

	for i := 0; i < b.N; i++ {
		userId, _ := gocloakClient.CreateUser(context.TODO(), AdminJWT.AccessToken, middleware.Realm, goCloakUser)
		_ = gocloakClient.DeleteUser(context.TODO(), AdminJWT.AccessToken, middleware.Realm, userId)
	}
}
