package routes

import (
	CO "Go-Mini-Social-Network/config"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Signup route
func Signup(c *gin.Context) {
	notLoggedIn(c)
	renderTemplate(c, "signup", gin.H{
		"title": "Signup For Free",
	})
}

// Login route
func Login(c *gin.Context) {
	notLoggedIn(c)
	renderTemplate(c, "login", gin.H{
		"title": "Login To Continue",
	})
}

// Logout route
func Logout(c *gin.Context) {
	loggedIn(c, "")
	session := CO.GetSession(c)
	delete(session.Values, "id")
	delete(session.Values, "username")
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/login")
}

// UserSignup function to register user
func UserSignup(c *gin.Context) {
	resp := make(map[string]interface{})

	username := strings.TrimSpace(c.PostForm("username"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := strings.TrimSpace(c.PostForm("password"))
	passwordAgain := strings.TrimSpace(c.PostForm("password_again"))

	mailErr := checkmail.ValidateFormat(email)

	db := CO.DB()

	var (
		userCount  int
		emailCount int
	)

	db.QueryRow("SELECT COUNT(id) AS userCount FROM users WHERE username=?", username).Scan(&userCount)
	db.QueryRow("SELECT COUNT(id) AS emailCount FROM users WHERE email=?", email).Scan(&emailCount)

	if username == "" || email == "" || password == "" || passwordAgain == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if len(username) < 4 || len(username) > 32 {
		resp["mssg"] = "Username should be between 4 and 32"
	} else if mailErr != nil {
		resp["mssg"] = "Invalid email format!!"
	} else if password != passwordAgain {
		resp["mssg"] = "Passwords don't match"
	} else if userCount > 0 {
		resp["mssg"] = "Username already exists!!"
	} else if emailCount > 0 {
		resp["mssg"] = "Email already exists!!"
	} else {

		stmt, _ := db.Prepare("INSERT INTO users(username, email, password, joined) VALUES (?, ?, ?, ?)")
		rs, iErr := stmt.Exec(username, email, hash(password), time.Now())
		CO.Err(iErr)
		insertID, _ := rs.LastInsertId()
		insStr := strconv.FormatInt(insertID, 10)

		dir, _ := os.Getwd()
		userPath := dir + "/public/users/" + insStr

		dirErr := os.Mkdir(userPath, 0655)
		CO.Err(dirErr)

		linkErr := os.Link(dir+"/public/images/golang.png", userPath+"/avatar.png")
		CO.Err(linkErr)

		session := CO.GetSession(c)
		session.Values["id"] = insStr
		session.Values["username"] = username
		session.Save(c.Request, c.Writer)

		resp["success"] = true
		resp["mssg"] = "Hello, " + username + "!!"

	}
	json(c, resp)
}

// UserLogin function to log user in
func UserLogin(c *gin.Context) {
	resp := make(map[string]interface{})

	rusername := strings.TrimSpace(c.PostForm("username"))
	rpassword := strings.TrimSpace(c.PostForm("password"))

	db := CO.DB()
	var (
		userCount int
		id        int
		username  string
		password  string
	)

	db.QueryRow("SELECT COUNT(id) AS userCount, id, username, password FROM users WHERE username=?", rusername).Scan(&userCount, &id, &username, &password)

	encErr := bcrypt.CompareHashAndPassword([]byte(password), []byte(rpassword))

	if rusername == "" || rpassword == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if userCount == 0 {
		resp["mssg"] = "Invalid username!!"
	} else if encErr != nil {
		resp["mssg"] = "Invalid password!!"
	} else {
		session := CO.GetSession(c)
		session.Values["id"] = strconv.Itoa(id)
		session.Values["username"] = username
		session.Save(c.Request, c.Writer)

		resp["mssg"] = "Hello, " + username + "!!"
		resp["success"] = true
	}
	json(c, resp)
}
