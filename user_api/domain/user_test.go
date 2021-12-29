package domain

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	passwordStr := "passwordHash"
	passwordHash := "$2a$10$FFopY1gpC5YV6R34UzBvY.HSoeucFn5irVG0iyEHAFLx65poGyslq"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserRepo := NewMockUserRepository(mockCtrl)
	testUser := UserService{
		repo: mockUserRepo,
	}

	mockUserRepo.EXPECT().Save("John", "Doe", "jDoe", gomock.Any()).Return(User{PasswordHash: passwordHash}, nil).Times(1)

	user, err := testUser.CreateUser("John", "Doe", "jDoe", passwordStr)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, user.IsCorrectPassword(passwordStr))

}
