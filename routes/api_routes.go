package routes

import (
	CO "Go-Mini-Social-Network/config"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
)

// CreateNewPost route
func CreateNewPost(c *gin.Context) {

	title := strings.TrimSpace(c.PostForm("title"))
	content := strings.TrimSpace(c.PostForm("content"))
	id, _ := CO.AllSessions(c)

	db := CO.DB()

	stmt, _ := db.Prepare("INSERT INTO posts(title, content, createdBy, createdAt) VALUES (?, ?, ?, ?)")
	rs, iErr := stmt.Exec(title, content, id, time.Now())
	CO.Err(iErr)

	insertID, _ := rs.LastInsertId()

	resp := map[string]interface{}{
		"postID": insertID,
		"mssg":   "Post Created!!",
	}
	json(c, resp)
}

// DeletePost route
func DeletePost(c *gin.Context) {
	post := c.PostForm("post")
	db := CO.DB()

	_, dErr := db.Exec("DELETE FROM posts WHERE postID=?", post)
	CO.Err(dErr)

	json(c, map[string]interface{}{
		"mssg": "Post Deleted!!",
	})
}

// UpdatePost route
func UpdatePost(c *gin.Context) {
	postID := c.PostForm("postID")
	title := c.PostForm("title")
	content := c.PostForm("content")

	db := CO.DB()
	db.Exec("UPDATE posts SET title=?, content=? WHERE postID=?", title, content, postID)

	json(c, map[string]interface{}{
		"mssg": "Post Updated!!",
	})
}

// UpdateProfile route
func UpdateProfile(c *gin.Context) {
	resp := make(map[string]interface{})

	id, _ := CO.AllSessions(c)
	username := strings.TrimSpace(c.PostForm("username"))
	email := strings.TrimSpace(c.PostForm("email"))
	bio := strings.TrimSpace(c.PostForm("bio"))

	mailErr := checkmail.ValidateFormat(email)
	db := CO.DB()

	if username == "" || email == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if mailErr != nil {
		resp["mssg"] = "Invalid email format!!"
	} else {
		_, iErr := db.Exec("UPDATE users SET username=?, email=?, bio=? WHERE id=?", username, email, bio, id)
		CO.Err(iErr)

		session := CO.GetSession(c)
		session.Values["username"] = username
		session.Save(c.Request, c.Writer)

		resp["mssg"] = "Profile updated!!"
		resp["success"] = true
	}

	json(c, resp)
}

// ChangeAvatar route
func ChangeAvatar(c *gin.Context) {
	resp := make(map[string]interface{})
	id, _ := CO.AllSessions(c)

	dir, _ := os.Getwd()
	dest := dir + "/public/users/" + id.(string) + "/avatar.png"

	dErr := os.Remove(dest)
	CO.Err(dErr)

	file, _ := c.FormFile("avatar")
	upErr := c.SaveUploadedFile(file, dest)

	if upErr != nil {
		resp["mssg"] = "An error occured!!"
	} else {
		resp["mssg"] = "Avatar changed!!"
		resp["success"] = true
	}

	json(c, resp)
}

// Follow route
func Follow(c *gin.Context) {
	id, _ := CO.AllSessions(c)
	user := c.PostForm("user")
	username := CO.Get(user, "username")

	db := CO.DB()
	stmt, _ := db.Prepare("INSERT INTO follow(followBy, followTo, followTime) VALUES(?, ?, ?)")
	_, exErr := stmt.Exec(id, user, time.Now())
	CO.Err(exErr)

	json(c, gin.H{
		"mssg": "Followed " + username + "!!",
	})
}

// Unfollow route
func Unfollow(c *gin.Context) {
	id, _ := CO.AllSessions(c)
	user := c.PostForm("user")
	username := CO.Get(user, "username")

	db := CO.DB()
	stmt, _ := db.Prepare("DELETE FROM follow WHERE followBy=? AND followTo=?")
	_, dErr := stmt.Exec(id, user)
	CO.Err(dErr)

	json(c, gin.H{
		"mssg": "Unfollowed " + username + "!!",
	})
}

// Like post route
func Like(c *gin.Context) {
	post := c.PostForm("post")
	db := CO.DB()
	id, _ := CO.AllSessions(c)

	stmt, _ := db.Prepare("INSERT INTO likes(postID, likeBy, likeTime) VALUES (?, ?, ?)")
	_, err := stmt.Exec(post, id, time.Now())
	CO.Err(err)

	json(c, gin.H{
		"mssg": "Post Liked!!",
	})
}

// Unlike post route
func Unlike(c *gin.Context) {
	post := c.PostForm("post")
	id, _ := CO.AllSessions(c)
	db := CO.DB()

	stmt, _ := db.Prepare("DELETE FROM likes WHERE postID=? AND likeBy=?")
	_, err := stmt.Exec(post, id)
	CO.Err(err)

	json(c, gin.H{
		"mssg": "Post Unliked!!",
	})
}
