package serverrequests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	. "goLang-coding-exercise-Final/consts"
	. "goLang-coding-exercise-Final/serverrrequests"
	. "goLang-coding-exercise-Final/structers"

	"github.com/stretchr/testify/assert"
)

var (
	users = map[int]User{
		1: {ID: 1, Name: "John Doe", Email: "johnDoe@example.com", Phone: "0505055", Address: "NoWhere"},
	}
	mutex = &sync.Mutex{}
)

func TestHandleSingleUser(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/users/John Doe", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleSingleUser(w, r, users, mutex)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	expected := `{"id":1,"name":"John Doe","email":"johnDoe@example.com","phone":"0505055","address":"NoWhere"}`
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")
}

func TestHandleSingleUser_UserNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/users/NonExistingUser", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleSingleUser(w, r, users, mutex)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code, "handler returned wrong status code")
}

func TestGetUserHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/users/John Doe", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetUserHandler(w, r, users)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	expected := `{"id":1,"name":"John Doe","email":"johnDoe@example.com","phone":"0505055","address":"NoWhere"}`
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")
}

func TestGetUserHandler_UserNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/users/NonExistingUser", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	users := make(map[int]User) // Initialize the users map
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetUserHandler(w, r, users)
	})
	fmt.Println(users)

	handler.ServeHTTP(rr, req)

	fmt.Println(users)
	assert.Equal(t, http.StatusNoContent, rr.Code, "handler returned wrong status code")
}

func TestDeleteUserHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/users/John Doe", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		DeleteUserHandler(w, r, users, mutex)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	_, exists := users[1]
	assert.False(t, exists, "handler failed to delete user")
}

func TestDeleteUserHandler_UserNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/users/NonExistingUser", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { DeleteUserHandler(w, r, users, mutex) })

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
}

func TestCreateUserHandler(t *testing.T) {
	newUser := User{ID: 2, Name: "Jane Doe", Email: "janeDoe@example.com", Phone: "0606066", Address: "SomeWhere"}
	body, err := json.Marshal(newUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(w, r, USER_PAGE, users, mutex)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")

	user, exists := users[2]
	assert.True(t, exists, "handler failed to create user")
	assert.Equal(t, "Jane Doe", user.Name, "handler failed to create user correctly")
}

func TestCreateUserHandler_InvalidData(t *testing.T) {
	newUser := User{Name: "", Email: "invalidEmail", Phone: "123", Address: ""}
	body, err := json.Marshal(newUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(w, r, USER_PAGE, users, mutex)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
}

func TestChangeUserHandler(t *testing.T) {
	changeUser := User{ID: 1, Name: "Johnny Doe", Email: "johnDoe@example.com", Phone: "0505055", Address: "NoWhere"}
	body, err := json.Marshal(changeUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/John Doe", bytes.NewBuffer(body))

	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ID", "1")
	req.Header.Set("Name", "somethingNotRelevant")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { ChangeUserHandler(w, r, users, mutex) })

	handler.ServeHTTP(rr, req)

	assert.NotEqual(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	user, exists := users[0]
	assert.False(t, exists, "handler failed to change user")
	assert.NotEqual(t, "Johnny Doe", user.Name, "handler failed to change user correctly")
}

func TestChangeUserHandler_UserNotFound(t *testing.T) {
	changeUser := User{ID: 99, Name: "NonExistingUser", Email: "nonexist@example.com", Phone: "0000000", Address: "NoWhere"}
	body, err := json.Marshal(changeUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/NonExistingUser", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("ID", "1")
	req.Header.Set("Name", "somethingNotRelevant")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ChangeUserHandler(w, r, users, mutex)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
}

func TestChangeUserHandler_InvalidData(t *testing.T) {
	changeUser := User{ID: 1, Name: "", Email: "invalidEmail", Phone: "123", Address: ""}
	body, err := json.Marshal(changeUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users/John Doe", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ID", "1")
	req.Header.Set("Name", "somethingNotRelevant")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ChangeUserHandler(w, r, users, mutex)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
}

func TestHandleSingleUser_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/users/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleSingleUser(w, r, nil, nil)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "handler returned wrong status code")
}

func TestHandleSingleUser_BadRequest(t *testing.T) {
	changeUser := User{ID: 1, Name: "", Email: "invalidEmail", Phone: "123", Address: ""}
	body, err := json.Marshal(changeUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	users := make(map[int]User)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleSingleUser(w, r, users, &sync.Mutex{})
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
}
