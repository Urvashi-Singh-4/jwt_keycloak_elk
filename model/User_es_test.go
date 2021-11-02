package model

import (
	"encoding/json"
	"testing"
)

func BenchmarkGetAllUsers(b *testing.B) {

	users := make([]User, 0)
	ElasticSearchClient()
	for i := 0; i < b.N; i++ {
		GetAllUsers(&users)
	}

}

func BenchmarkCreateandDeleteUsers(b *testing.B) {
	ElasticSearchClient()
	dataByte, _ := json.Marshal(map[string]string{"hello": "world"})

	for i := 0; i < b.N; i++ {
		user, _ := CreateNewUser(string(dataByte))
		DeleteUserByID(user.Id)

	}

}
