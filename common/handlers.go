package common

import(
	"database/sql"
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

type UserInfo struct {
	Name string `json:"username"`
	Email string `json:"email"`
}

func LoadFile(fileName string) (string, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func LoginPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("123")
	var body, _ = LoadFile("templates/login.html")
	fmt.Fprintf(response, body)
}

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username := request.FormValue("username")
	password := request.FormValue("password")

	if len(username) <= 0 || len(password) <= 0 {
		http.Error(response, "These fields cannot be blank!", 500)
		return
	}

	var dbUsername string
	var dbPassword string
	err := db.QueryRow("SELECT Username, Password FROM Users WHERE Username=?", username).Scan(&dbUsername, &dbPassword)
	if err != nil {
		http.Redirect(response, request, "/login?retry=1", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if (err != nil) {
		http.Error(response, "Password is not matched", 301)
		return
	}

	// Redirect to "index.html" page
	SetCookie(username, response)
	http.Redirect(response, request, "/search", 301)
}

func SignupPageHandler(response http.ResponseWriter, request *http.Request) {
	var body, _ = LoadFile("templates/signup.html")
	fmt.Fprintf(response, body)
}

func SignupHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username := request.FormValue("username")
	password := request.FormValue("password")
	confirmPassword := request.FormValue("confirm-password")
	email := request.FormValue("email")

	// Simple validate.
	if len(username) <= 0 || len(password) <= 0 || len(confirmPassword) <=0 || len(email) <=0 {
		http.Error(response, "These fields cannot be blank!", 500)
		return
	}

	// Compare password and confirm password.
	if password != confirmPassword {
		http.Error(response, "Password and confirm password are not matched.", 500)
		return
	}

	// Find username in DB.
	var user string
	err := db.QueryRow("SELECT Username FROM Users WHERE Username=?", username).Scan(&user)
	log.Println(err)
	switch {
	case err == sql.ErrNoRows:
		// If not found, 
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(response, "Hash error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO Users(Username, Password, Email) VALUES(?, ?, ?)", username, hashedPassword, email)
		if err != nil {
			http.Error(response, "Insert error, unable to create your account.", 500)
			return
		}

		response.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(response, "Existing user error, unable to create your account.", 500)
		return
	default:
		http.Redirect(response, request, "/", 301)
	}
}

func SearchPageHandler(response http.ResponseWriter, request *http.Request) {
	username := GetUsername(request)

	if len(username) > 0 {
		var body, _ = LoadFile("templates/search.html")
		fmt.Fprintf(response, body)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func SearchHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username := request.FormValue("username")

	var rows *sql.Rows
	var err error
	if len(username) <= 0 {
		rows, err = db.Query("SELECT Username, Email FROM Users")
		if err != nil {
			log.Panic(err)
			http.Error(response, "Cannot query database.", 500)
			return
		}
	} else {
		// Query username from DB.
		username = "%" + username + "%"
		rows, err = db.Query("SELECT Username, Email FROM Users WHERE Username LIKE ?", username)
		if err != nil {
			log.Panic(err)
			http.Error(response, "Cannot query database.", 500)
			return
		}
	}

	// Extract data from query result.
	defer rows.Close()
	data := make([]UserInfo, 0)
	for rows.Next() {
		var userInfo UserInfo
		err = rows.Scan(&userInfo.Name, &userInfo.Email)
		if err != nil {
			log.Panic(err)
		} else {
			data = append(data, userInfo)
		}
	}

	// Convert data to json.
	json, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
		return
	}

	response.Write(json)
	// fmt.Fprintf(response, "Result: %s", json)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	ClearCookie(response)
	http.Redirect(response, request, "/", 302)
}

func SetCookie(username string, response http.ResponseWriter) {
	value := map[string]string {
		"name": username,
	}

	if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie {
			Name: "cookie",
			Value: encoded,
			Path: "/",
		}
		http.SetCookie(response, cookie)
	}
}

func ClearCookie(response http.ResponseWriter) {
	cookie := &http.Cookie {
		Name: "cookie",
		Value: "",
		Path: "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func GetUsername(request *http.Request) (username string) {
	if cookie, err := request.Cookie("cookie"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
		}
	}

	return username
}