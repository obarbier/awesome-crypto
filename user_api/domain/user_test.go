package domain

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserRepo := NewMockUserRepository(mockCtrl)
	testUser := UserService{
		repo: mockUserRepo,
	}

	ctx := context.Background()
	mockUserRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	user, err := testUser.CreateUser(ctx, "John", "Doe", "jDoe", "passwordStr")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, user.IsCorrectPassword("passwordStr"))

}
