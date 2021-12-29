package adapters

import (
	"context"
	"testing"
)

var repo, _ = NewMongoRepository()

func TestCreateUser(t *testing.T) {
	_, err := repo.Save(context.Background(), "John", "Doe", "jDoe", "passwordHash")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserByUserId(t *testing.T) {
	_, err := repo.Get(context.Background(), "df3f55b4-6853-11ec-84e8-c85b768b87cd")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateUserByUserId(t *testing.T) {
	_, err := repo.Update(context.Background(), "df3f55b4-6853-11ec-84e8-c85b768b87cd", "John2", "Doe2", "jDoe2", "passwordHash2")
	if err != nil {
		t.Fatal(err)
	}
}
