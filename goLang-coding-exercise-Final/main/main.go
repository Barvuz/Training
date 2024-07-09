/*
 _____  ____  _                                   _ _                                      _
/ ____|/ __ \| |                                 | (_)                                    (_)
| |  __| |  | | | __ _ _ __   __ _    ___ ___   __| |_ _ __   __ _    _____  _____ _ __ ___ _ ___  ___
| | |_ | |  | | |/ _` | '_ \ / _` |  / __/ _ \ / _` | | '_ \ / _` |  / _ \ \/ / _ \ '__/ __| / __|/ _ \
| |__| | |__| | | (_| | | | | (_| | | (_| (_) | (_| | | | | | (_| | |  __/>  <  __/ | | (__| \__ \  __/
\_____|\____/|_|\__,_|_| |_|\__, |  \___\___/ \__,_|_|_| |_|\__, |  \___/_/\_\___|_|  \___|_|___/\___|
							 __/ |                           __/ |
							|___/                           |___/
*/

package main

import (
	"fmt"
	. "go-user-server/consts"
	dbf "go-user-server/dataBase"
	sq "go-user-server/serverRequests"
	. "go-user-server/structers"
	"log"
	"net/http"
	"sync"
)

var (
	users        = make(map[int]User)
	mutex_helper sync.Mutex
)

/*
Function: main
Input: None
Output: None
*/
func main() {
	fmt.Println("WORKING")

	users, err := dbf.LoadUsers(USER_PAGE) //load all the users from the file "users.json"
	if err != nil {
		fmt.Println("Error loading users:", err)
		return
	}

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) { sq.HandleUsers(w, r, users, &mutex_helper) })
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) { sq.HandleSingleUser(w, r, users, &mutex_helper) })

	fmt.Println("Do you want to change any user? 1 for yes 0 for no")

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
