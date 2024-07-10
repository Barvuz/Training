package databasefunctions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	. "goLang-coding-exercise-Final/structers"

	"github.com/go-playground/validator/v10"
)

func ReadUsersFromFile(filename string) (map[int]User, error) { //NOT IN USE
	// Open the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Decode JSON into a slice of User structs
	var users []User
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	// Convert slice to map[string]User
	userMap := make(map[int]User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	return userMap, nil
}

/*
Function: getAllUsersHandler
Input: w (http.ResponseWriter), _ (*http.Request)
Output: None
*/
func GetAllUsersHandler(w http.ResponseWriter, _ *http.Request, users map[int]User) {
	var names []string
	for _, user := range users {
		names = append(names, user.Name) //add every user to names array
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names) //write on the screen the name (in a json format to be more readable)
}

/*
Function: loadUsers
Input: filename (string)
Output: error
*/
func LoadUsers(filename string) (map[int]User, error) {
	users, err := ReadUsersFromFile(filename)
	if err != nil {
		return nil, errors.New("error in reading users from file")
	}
	fmt.Println(users)

	seenIDs := make(map[int]bool) // Handle duplicates
	for _, user := range users {
		if seenIDs[user.ID] {
			fmt.Println("duplicate ID found: ", user.ID)
			return nil, errors.New("duplicate ID found")
		}
		seenIDs[user.ID] = true
		// Initialize the validator
		validate := validator.New()
		// Example user
		userValid := user // No need to take a pointer
		// Validate the struct
		errVal := validate.Struct(userValid)
		if errVal != nil {
			// Handling validation errors
			for _, err := range errVal.(validator.ValidationErrors) {
				fmt.Printf("Field '%s' is invalid: %s\n", err.Field(), err.ActualTag())
			}
			return nil, errors.New("invalid field")
		}
	}

	return users, nil
}

/*
Function: saveUsers
Input: filename (string)
Output: error
*/
func SaveUsers(filename string, mutex_helper *sync.Mutex, users map[int]User) error {
	mutex := mutex_helper
	mutex.Lock() //rotect the users map while marshaling user data and writing it to the file.
	defer mutex.Unlock()
	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}

	data, err := json.Marshal(userList) //encodes the user struct into a JSON-encoded byte slice.
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644) //save the users in user.json file (0644 = read and write mode)
}
