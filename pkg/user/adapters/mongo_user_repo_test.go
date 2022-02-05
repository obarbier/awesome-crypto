package adapters

import (
	"context"
	"github.com/obarbier/awesome-crypto/pkg/user/domain"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var repo, _ = NewMongoRepository()

func TestUserManagement(t *testing.T) {
	user := domain.User{
		Id:           "unit-test-id",
		FirstName:    "John",
		LastName:     "Doe",
		UserId:       "jDoe",
		PasswordHash: "passwordHash",
	}
	err := repo.Save(context.Background(), &user)
	if err != nil {
		t.Fatal(err)
	}
	id := user.Id

	userUpdate := &domain.User{
		UserId: "jDoeUpdate",
	}
	err = repo.UpdateByID(context.Background(), id, userUpdate)
	if err != nil {
		t.Fatal(err)
	}

	user2, err := repo.FindById(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, user.FirstName, user2.FirstName)
	assert.Equal(t, user.LastName, user2.LastName)
	assert.NotEqual(t, user.UserId, user2.UserId)
	assert.Equal(t, userUpdate.UserId, user2.UserId)
	assert.Equal(t, user.PasswordHash, user2.PasswordHash)

	err = repo.DeleteById(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	emptyUser, err := repo.FindById(context.Background(), id)
	if err != nil {
		assert.True(t, strings.Contains(err.Error(), "no documents in result"))
	}
	assert.Nil(t, emptyUser, "user should not be in db")

}

func TestUserManagementInvalidFlow(t *testing.T) {
	user := domain.User{
		Id:           "unit-test-id",
		FirstName:    "John",
		LastName:     "Doe",
		UserId:       "jDoe",
		PasswordHash: "passwordHash",
	}
	err := repo.Save(context.Background(), &user)
	if err != nil {
		t.Fatal(err)
	}
	id := user.Id
	userUpdate := &domain.User{
		UserId: "jDoeUpdate",
	}
	err = repo.UpdateByID(context.Background(), "Not-valid-Id", userUpdate)
	if err != nil {
		assert.True(t, strings.Contains(err.Error(), "no user with Id"))
	}

	err = repo.DeleteById(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

}
