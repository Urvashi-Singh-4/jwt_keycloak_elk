package model

import (
	"context"
	"encoding/json"
	"errors"

	elastic "github.com/olivere/elastic/v7"
)

type User struct {
	Id       string `json:"Id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"_"`
}

func (u *User) IndexName() string {
	return "user"
}

func CreateNewUser(dataJson string) (*elastic.IndexResponse, error) {
	var u User
	x, err := DB.Index().Index(u.IndexName()).BodyJson(dataJson).Do(context.Background())
	if err != nil {
		return x, err
	}

	return x, nil
}

func GetAllUsers(users *[]User) error {
	var u User

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchAllQuery())

	searchResult, err := DB.Search().Index(u.IndexName()).SearchSource(searchSource).Do(context.TODO())
	if err != nil {
		return err
	}

	for _, hit := range searchResult.Hits.Hits {
		var user User

		err = json.Unmarshal(hit.Source, &user)
		if err != nil {
			return err
		}

		*users = append(*users, user)
	}

	return nil
}

func DeleteUserByID(userID string) (error, string) {
	var u User
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("id", userID))
	searchResult, err := DB.Search().Index(u.IndexName()).SearchSource(searchSource).Do(context.Background())
	if err != nil {
		return err, "Error occured while deleting the user"
	}
	if len(searchResult.Hits.Hits) == 0 {
		return errors.New("USER_NOT_FOUND"), "user not found"
	}
	for _, hit := range searchResult.Hits.Hits {
		_, err = DB.Delete().Index(u.IndexName()).Id(hit.Id).Do(context.Background())
		if err != nil {
			return err, "Error occured while deleting the user"
		}
	}
	return err, "User deleted successfully"
}

func UpdateUserByID(userID string, updateData map[string]interface{}) (error, string) {
	var u User
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("id", userID))
	searchResult, err := DB.Search().Index(u.IndexName()).SearchSource(searchSource).Do(context.Background())
	if err != nil {
		return err, "Error occured while updating the user"
	}
	if len(searchResult.Hits.Hits) == 0 {
		return errors.New("USER_NOT_FOUND"), "User not found"
	}
	for _, hit := range searchResult.Hits.Hits {
		_, err = DB.Update().Index(u.IndexName()).Id(hit.Id).Doc(updateData).Do(context.Background())
		if err != nil {
			return err, "Error occured while updating the user"
		}
	}
	return err, "User updated successfully"
}
