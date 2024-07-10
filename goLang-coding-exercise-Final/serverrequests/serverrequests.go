package serverrequests

import (
	"encoding/json"
	"errors"
	"fmt"
	. "goLang-coding-exercise-Final/consts"
	dbf "goLang-coding-exercise-Final/database"
	. "goLang-coding-exercise-Final/structers"
	. "goLang-coding-exercise-Final/usermethod"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
)

/*
Function: getUserHandler
Input: w (http.ResponseWriter), r (*http.Request)
Output: None
*/
func GetUserHandler(w http.ResponseWriter, r *http.Request, users map[int]User) error {
	name := r.URL.Path[len("/users/"):]
	user := IsNameExist(users, name)
	if user == nil {
		http.Error(w, "user not exist", http.StatusNoContent)
		return errors.New("user not exist")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(*user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, users map[int]User, mutex_helper *sync.Mutex) error {
	name := r.URL.Path[len("/users/"):]

	user := IsNameExist(users, name)
	if user == nil {
		http.Error(w, "user not exist", http.StatusNotFound)
		return errors.New("user not exist")
	}

	mutex_helper.Lock()
	delete(users, user.ID)
	mutex_helper.Unlock()
	dbf.SaveUsers(USER_PAGE, mutex_helper, users)
	return nil
}

func HandleSingleUser(w http.ResponseWriter, r *http.Request, users map[int]User, mutexHelper *sync.Mutex) error {
	var err error

	switch r.Method {
	case http.MethodGet:
		err = GetUserHandler(w, r, users)
	case http.MethodPost:
		err = CreateUserHandler(w, r, USER_PAGE, users, mutexHelper)
	case http.MethodDelete:
		err = DeleteUserHandler(w, r, users, mutexHelper)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		err = errors.New("invalid request method (unknown)")
	}

	if err != nil {
		if err.Error() == "user not exist" {
			http.Error(w, fmt.Sprintf("Bad Request (L): %s", err.Error()), http.StatusNoContent)
			return errors.New("user not exist")
		}
		fmt.Printf("Something not working %s\n", err)
		http.Error(w, fmt.Sprintf("Bad Request (L): %s", err.Error()), http.StatusBadRequest)
		return errors.New("Bad Request (L) Something not working \n")
	}

	return nil
}

/*
Function: createUserHandler
Input: w (http.ResponseWriter), r (*http.Request)
Output: None
*/
func CreateUserHandler(w http.ResponseWriter, r *http.Request, fileName string, users map[int]User, mutex_helper *sync.Mutex) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return errors.New("invalid request method")
	}

	// Decode JSON request body into User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return errors.New("invalid request body")
	}

	// Initialize the validator
	validate := validator.New()

	// Validate the struct
	errVal := validate.Struct(&user)
	if errVal != nil {
		// Handling validation errors
		for _, errVal := range errVal.(validator.ValidationErrors) {
			fmt.Printf("Field '%s' is invalid: %s\n", errVal.Field(), errVal.ActualTag())
		}
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return errors.New("invalid user data")
	}

	fmt.Println("User is valid!")
	mutex_helper.Lock()
	users[user.ID] = user
	mutex_helper.Unlock()

	err = dbf.SaveUsers(fileName, mutex_helper, users)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return errors.New("error saving user")
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func ChangeUserHandler(w http.ResponseWriter, r *http.Request, users map[int]User, mutex_helper *sync.Mutex) {
	if r.Header.Get("ID") == "" {
		http.Error(w, "Missing ID param", http.StatusMethodNotAllowed)
	}
	if r.Header.Get("Name") == "" {
		http.Error(w, "Missing Name param", http.StatusMethodNotAllowed)
	}

	id, err := strconv.Atoi(r.Header.Get("ID")) //change to url.query/ header
	if err != nil {
		// Handle the error if strconv.Atoi fails to convert idStr to an integer
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	user, ok := users[id]
	if !ok {
		// Handle case where user with given ID does not exist in the map
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}

	user.Name = r.Header.Get("Name")
	users[id] = user

	dbf.SaveUsers(USER_PAGE, mutex_helper, users)

}

func HandleUsers(w http.ResponseWriter, r *http.Request, users map[int]User, mutexHelper *sync.Mutex) {
	var err error = nil

	switch r.Method {
	case http.MethodGet:
		dbf.GetAllUsersHandler(w, r, users)
	case http.MethodPost:
		if r.URL.Query().Get("Command") == "insert" {
			err = CreateUserHandler(w, r, USER_PAGE, users, mutexHelper)
		} else if r.URL.Query().Get("Command") == "change" {
			ChangeUserHandler(w, r, users, mutexHelper)
		} else {
			fmt.Println(r.URL.Query().Get("Command"))
			http.Error(w, "Missing Command param", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		err = errors.New("invalid request method (unknown)")
	}

	if err != nil {
		fmt.Printf("Something not working %s\n", err)
		http.Error(w, fmt.Sprintf("Bad Request (L): %s", err.Error()), http.StatusBadRequest)
		return
	}

}
