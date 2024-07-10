package databasefunctions_test

import (
	. "goLang-coding-exercise-Final/database"
	. "goLang-coding-exercise-Final/structers"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWriteUsers(t *testing.T) {
	filename := "test_users.json"
	expectedUsers := map[int]User{
		1: {ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "123456789", Address: "123 Street"},
		2: {ID: 2, Name: "Bob", Email: "bob@example.com", Phone: "987654321", Address: "456 Avenue"},
	}

	// Test saving users to file
	mutex := sync.Mutex{}
	err := SaveUsers(filename, &mutex, expectedUsers)
	assert.NoError(t, err, "SaveUsers should not return an error")

	// Test loading users from file
	actualUsers, err := LoadUsers(filename)
	assert.NoError(t, err, "LoadUsers should not return an error")
	assert.Equal(t, len(expectedUsers), len(actualUsers), "Loaded users count should match expected")

	for id, expectedUser := range expectedUsers {
		actualUser, ok := actualUsers[id]
		assert.True(t, ok, "Loaded user with ID %d should exist", id)
		assert.Equal(t, expectedUser.ID, actualUser.ID, "IDs should match")
		assert.Equal(t, expectedUser.Name, actualUser.Name, "Names should match")
		assert.Equal(t, expectedUser.Email, actualUser.Email, "Emails should match")
		assert.Equal(t, expectedUser.Phone, actualUser.Phone, "Phones should match")
		assert.Equal(t, expectedUser.Address, actualUser.Address, "Addresses should match")
	}

	// Additional test cases
	t.Run("Test with missing Name", func(t *testing.T) {
		missingNameUsers := map[int]User{
			1: {ID: 1, Email: "alice@example.com", Phone: "123456789", Address: "123 Street"},
			2: {ID: 2, Email: "bob@example.com", Phone: "987654321", Address: "456 Avenue"},
		}

		err := SaveUsers(filename, &mutex, missingNameUsers)
		assert.NoError(t, err, "SaveUsers should not return an error")

		actualUsers, err := LoadUsers(filename)
		assert.Error(t, err, "LoadUsers should return an error for missing Name")
		assert.Nil(t, actualUsers, "Actual users should be nil for invalid data")
	})

	t.Run("Test with missing Email", func(t *testing.T) {
		missingEmailUsers := map[int]User{
			1: {ID: 1, Name: "Alice", Phone: "123456789", Address: "123 Street"},
			2: {ID: 2, Name: "Bob", Phone: "987654321", Address: "456 Avenue"},
		}

		err := SaveUsers(filename, &mutex, missingEmailUsers)
		assert.NoError(t, err, "SaveUsers should not return an error")

		actualUsers, err := LoadUsers(filename)
		assert.Error(t, err, "LoadUsers should return an error for missing Email")
		assert.Nil(t, actualUsers, "Actual users should be nil for invalid data")
	})

	t.Run("Test with missing Phone", func(t *testing.T) {
		missingPhoneUsers := map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com", Address: "123 Street"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com", Address: "456 Avenue"},
		}

		err := SaveUsers(filename, &mutex, missingPhoneUsers)
		assert.NoError(t, err, "SaveUsers should not return an error")

		actualUsers, err := LoadUsers(filename)
		assert.Error(t, err, "LoadUsers should return an error for missing Phone")
		assert.Nil(t, actualUsers, "Actual users should be nil for invalid data")
	})

	t.Run("Test with missing Address", func(t *testing.T) {
		missingAddressUsers := map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "123456789"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com", Phone: "987654321"},
		}

		err := SaveUsers(filename, &mutex, missingAddressUsers)
		assert.NoError(t, err, "SaveUsers should not return an error")

		actualUsers, err := LoadUsers(filename)
		assert.Error(t, err, "LoadUsers should return an error for missing Address")
		assert.Nil(t, actualUsers, "Actual users should be nil for invalid data")
	})

	// Clean up: Remove the test file
	err = os.Remove(filename)
	assert.NoError(t, err, "Error deleting test file")
}

func TestReadUsersFromFile(t *testing.T) {
	filename := "test_read_users.json"
	expectedUsers := map[int]User{
		1: {ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "123456789", Address: "123 Street"},
		2: {ID: 2, Name: "Bob", Email: "bob@example.com", Phone: "987654321", Address: "456 Avenue"},
	}

	// Save expected users to file
	err := SaveUsers(filename, &sync.Mutex{}, expectedUsers)
	assert.NoError(t, err, "SaveUsers should not return an error")

	// Test reading users from file
	actualUsers, err := ReadUsersFromFile(filename)
	assert.NoError(t, err, "ReadUsersFromFile should not return an error")
	assert.Equal(t, len(expectedUsers), len(actualUsers), "Read users count should match expected")

	for id, expectedUser := range expectedUsers {
		actualUser, ok := actualUsers[id]
		assert.True(t, ok, "Read user with ID %d should exist", id)
		assert.Equal(t, expectedUser.ID, actualUser.ID, "IDs should match")
		assert.Equal(t, expectedUser.Name, actualUser.Name, "Names should match")
		assert.Equal(t, expectedUser.Email, actualUser.Email, "Emails should match")
		assert.Equal(t, expectedUser.Phone, actualUser.Phone, "Phones should match")
		assert.Equal(t, expectedUser.Address, actualUser.Address, "Addresses should match")
	}

	// Clean up: Remove the test file
	err = os.Remove(filename)
	assert.NoError(t, err, "Error deleting test file")
}
