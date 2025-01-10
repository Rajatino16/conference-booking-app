package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupUserRepository() Repository {
	return NewInMemoryRepository()
}

func setupUserService() Service {
	repo := setupUserRepository()
	return NewService(repo)
}

func TestAddUser(t *testing.T) {
	service := setupUserService()

	// Add a user
	err := service.AddUser(AddUserRequest{
		ID: "user1",
	})

	// Verify the user is added successfully
	assert.NoError(t, err)
}

func TestAddDuplicateUser(t *testing.T) {
	service := setupUserService()

	// Add a user
	err := service.AddUser(AddUserRequest{
		ID: "user1",
	})
	assert.NoError(t, err)

	// Attempt to add another user with the same ID
	err = service.AddUser(AddUserRequest{
		ID: "user1",
	})

	// Verify that the second addition fails with a conflict error
	assert.Error(t, err)
	assert.Equal(t, "resource conflict", err.Error())
}
