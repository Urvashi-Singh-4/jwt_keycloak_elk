package middleware

import (
	"context"
	"strings"

	"github.com/Nerzal/gocloak/v8"
)

const Realm = "master"
const KeyCloakClient = "elastic-go"
const KeyCloakClientSecret = "526fa944-d711-4531-9755-47c585f17436"

func CheckAuthorised(accesstoken string, requiredPermission string) bool {

	gocloakClient := gocloak.NewClient("http://localhost:8080")
	_, claims, err := gocloakClient.DecodeAccessToken(context.Background(), accesstoken, Realm, KeyCloakClient)
	if err != nil {

		return false
	}

	jwt, err := gocloakClient.Login(context.Background(),
		KeyCloakClient,
		KeyCloakClientSecret,
		Realm,
		"admin",
		"admin",
	)
	role := (*claims)["role"].([]interface{})[0].(string)

	r, err := gocloakClient.GetRealmRole(context.Background(), jwt.AccessToken, Realm, role)
	if err != nil {
		return false
	}

	for _, p := range strings.Split((*r.Attributes)["permission"][0], ",") {
		if p == requiredPermission {
			return true
		}
	}

	return false
}

func AdminLogin() (*gocloak.JWT, error) {
	gocloakClient := gocloak.NewClient("http://localhost:8080")
	return gocloakClient.Login(context.Background(),
		KeyCloakClient,
		KeyCloakClientSecret,
		Realm,
		"admin",
		"admin",
	)

}
