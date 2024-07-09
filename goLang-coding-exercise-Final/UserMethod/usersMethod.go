package UserMethod

import (
	. "go-user-server/structers"
	"sort"
)

/*
Function: sortUsersListByID
Input: None
Output: []User
*/
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
