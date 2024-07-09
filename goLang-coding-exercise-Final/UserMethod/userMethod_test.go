package UserMethod_test

import (
	. "go-user-server/structers"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SortUsersListByID(users map[string]User) map[string]User {
	// Extract values from map to a slice
	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}

	// Define a sorting function
	sort.Slice(userList, func(i, j int) bool {
		return userList[i].ID < userList[j].ID
	})

	// Create a new sorted map
	sortedUsers := make(map[string]User)
	for _, user := range userList {
		sortedUsers[user.Name] = user
	}

	return sortedUsers
}

func IsNameExist(users map[int]User, name string) *User {
	for _, user := range users {
		if user.Name == name {
			return &user
		}
	}
	return nil
}

func TestSortUsersListByID(t *testing.T) {
	users := map[string]User{
		"Alice":   {ID: 3, Name: "Alice", Email: "alice@example.com", Phone: "123456789", Address: "123 Street"},
		"Bob":     {ID: 1, Name: "Bob", Email: "bob@example.com", Phone: "987654321", Address: "456 Avenue"},
		"Charlie": {ID: 2, Name: "Charlie", Email: "charlie@example.com", Phone: "123123123", Address: "789 Boulevard"},
	}

	sortedUsers := SortUsersListByID(users)
	expectedOrder := []string{"Bob", "Charlie", "Alice"}

	i := 0
	for _, user := range sortedUsers {
		assert.Equal(t, expectedOrder[i], user.Name, "User order should be sorted by ID")
		i++
	}
}

func TestIsNameExist(t *testing.T) {
	users := map[int]User{
		1: {ID: 1, Name: "Alice", Email: "alice@example.com", Phone: "123456789", Address: "123 Street"},
		2: {ID: 2, Name: "Bob", Email: "bob@example.com", Phone: "987654321", Address: "456 Avenue"},
	}

	t.Run("existing user", func(t *testing.T) {
		user := IsNameExist(users, "Alice")
		assert.NotNil(t, user, "Expected to find user Alice")
		assert.Equal(t, "Alice", user.Name, "User name should be Alice")
	})

	t.Run("non-existing user", func(t *testing.T) {
		user := IsNameExist(users, "Charlie")
		assert.Nil(t, user, "Expected not to find user Charlie")
	})

	t.Run("empty name", func(t *testing.T) {
		user := IsNameExist(users, "")
		assert.Nil(t, user, "Expected not to find an empty name")
	})

	t.Run("case sensitivity", func(t *testing.T) {
		user := IsNameExist(users, "alice")
		assert.Nil(t, user, "Expected not to find user alice with different case")
	})
}
