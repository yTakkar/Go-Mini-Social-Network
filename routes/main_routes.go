package routes

import (
	CO "Go-Mini-Social-Network/config"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Index route
func Index(c *gin.Context) {
	loggedIn(c, "/welcome")

	id, _ := CO.AllSessions(c)
	db := CO.DB()
	var (
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)
	feeds := []interface{}{}

	stmt, _ := db.Prepare("SELECT posts.postID, posts.title, posts.content, posts.createdBy, posts.createdAt from posts, follow WHERE follow.followBy=? AND follow.followTo = posts.createdBy ORDER BY posts.postID DESC")
	rows, qErr := stmt.Query(id)
	CO.Err(qErr)

	for rows.Next() {
		rows.Scan(&postID, &title, &content, &createdBy, &createdAt)
		feed := map[string]interface{}{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
		}
		feeds = append(feeds, feed)
	}

	renderTemplate(c, "index", gin.H{
		"title":   "Home",
		"session": ses(c),
		"posts":   feeds,
		"GET":     CO.Get,
	})
}

// Welcome route
func Welcome(c *gin.Context) {
	notLoggedIn(c)
	renderTemplate(c, "welcome", gin.H{
		"title": "Welcome",
	})
}

// NotFound route
func NotFound(c *gin.Context) {
	renderTemplate(c, "404", gin.H{
		"title":   "Oops!! Error",
		"session": ses(c),
	})
}

// Profile Page
func Profile(c *gin.Context) {
	loggedIn(c, "")

	user := c.Param("id")
	sesID, _ := CO.AllSessions(c)
	db := CO.DB()

	// VARS FOR USER DETAILS
	var (
		userCount int
		userID    int
		username  string
		email     string
		bio       string
	)

	// VARS FOR POSTS
	var (
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)
	posts := []interface{}{}

	var (
		followers  int //for followers
		followings int //for followings
		pViews     int // for profile views
	)

	me := CO.MeOrNot(c, user) // Check if its me or not
	var noMssg string         // Mssg to be displayed when user has no posts

	if me == true {
		noMssg = "You have no posts. Go ahead and create one!!"
	} else {
		noMssg = username + " has no posts!!"

		// VIEW PROFILE
		if sesID != nil {
			stmt, _ := db.Prepare("INSERT INTO profile_views(viewBy, viewTo, viewTime) VALUES(?, ?, ?)")
			_, pvErr := stmt.Exec(sesID, user, time.Now())
			CO.Err(pvErr)
		}

	}

	// USER DETAILS
	db.QueryRow("SELECT COUNT(id) AS userCount, id AS userID, username, email, bio FROM users WHERE id=?", user).Scan(&userCount, &userID, &username, &email, &bio)
	invalid(c, userCount)

	// POSTS
	stmt, sErr := db.Prepare("SELECT * FROM posts WHERE createdBy=? ORDER BY postID DESC")
	CO.Err(sErr)
	rows, gErr := stmt.Query(userID)
	CO.Err(gErr)

	for rows.Next() {
		rows.Scan(&postID, &title, &content, &createdBy, &createdAt)
		post := map[string]interface{}{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
		}
		posts = append(posts, post)
	}

	db.QueryRow("SELECT COUNT(followID) AS followers FROM follow WHERE followTo=?", user).Scan(&followers)  // FOLLOWERS
	db.QueryRow("SELECT COUNT(followID) AS followers FROM follow WHERE followBy=?", user).Scan(&followings) // FOLLOWINGS
	db.QueryRow("SELECT COUNT(viewID) AS pViews FROM profile_views WHERE viewTo=?", user).Scan(&pViews)     // PROFILE VIEWS

	renderTemplate(c, "profile", gin.H{
		"title":   "@" + username,
		"session": ses(c),
		"user": gin.H{
			"id":       strconv.Itoa(userID),
			"username": username,
			"email":    email,
			"bio":      bio,
		},
		"posts":      posts,
		"followers":  followers,
		"followings": followings,
		"views":      pViews,
		"no_mssg":    noMssg,
		"GET":        CO.Get,
		"isF":        CO.IsFollowing,
	})

}

// Explore route
func Explore(c *gin.Context) {
	loggedIn(c, "")
	user, _ := CO.AllSessions(c)
	db := CO.DB()
	var (
		id       int
		username string
		email    string
	)
	explore := []interface{}{}

	stmt, _ := db.Prepare("SELECT id, username, email FROM users WHERE id <> ? ORDER BY RAND() LIMIT 10")
	rows, err := stmt.Query(user)
	CO.Err(err)

	for rows.Next() {
		rows.Scan(&id, &username, &email)
		exp := map[string]interface{}{
			"id":       id,
			"username": username,
			"email":    email,
		}
		explore = append(explore, exp)
	}

	renderTemplate(c, "explore", gin.H{
		"title":   "Explore",
		"session": ses(c),
		"users":   explore,
		"GET":     CO.Get,
		"noF":     CO.NoOfFollowers,
		"UD":      CO.UsernameDecider,
	})
}

// CreatePost route
func CreatePost(c *gin.Context) {
	loggedIn(c, "")
	renderTemplate(c, "create_post", gin.H{
		"title":   "Create Post",
		"session": ses(c),
	})
}

// ViewPost route
func ViewPost(c *gin.Context) {
	loggedIn(c, "")

	param := c.Param("id")
	db := CO.DB()
	var (
		postCount int
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)
	var likesCount int

	// post details
	db.QueryRow("SELECT COUNT(postID) AS postCount, postID, title, content, createdBy, createdAt FROM posts WHERE postID=?", param).Scan(&postCount, &postID, &title, &content, &createdBy, &createdAt)
	invalid(c, postCount)

	// likes
	db.QueryRow("SELECT COUNT(likeID) AS likesCount FROM likes WHERE postID=?", param).Scan(&likesCount)

	renderTemplate(c, "view_post", gin.H{
		"title":   "View Post",
		"session": ses(c),
		"post": gin.H{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
		},
		"postCreatedBy": strconv.Itoa(createdBy),
		"lon":           CO.LikedOrNot,
		"likes":         likesCount,
	})
}

// EditPost route
func EditPost(c *gin.Context) {
	loggedIn(c, "")

	post := c.Param("id")
	db := CO.DB()
	var (
		postCount int
		postID    int
		title     string
		content   string
	)

	db.QueryRow("SELECT COUNT(postID) AS postCount, postID, title, content FROM posts WHERE postID=?", post).Scan(&postCount, &postID, &title, &content)
	invalid(c, postCount)

	renderTemplate(c, "edit_post", gin.H{
		"title":   "Edit Post",
		"session": ses(c),
		"post": gin.H{
			"postID":  postID,
			"title":   title,
			"content": content,
		},
	})
}

// EditProfile route
func EditProfile(c *gin.Context) {
	loggedIn(c, "")

	db := CO.DB()
	id, _ := CO.AllSessions(c)
	var (
		email  string
		bio    string
		joined string
	)
	db.QueryRow("SELECT email, bio, joined FROM users WHERE id=?", id).Scan(&email, &bio, &joined)
	renderTemplate(c, "edit_profile", gin.H{
		"title":   "Edit Profile",
		"session": ses(c),
		"email":   email,
		"bio":     bio,
		"joined":  joined,
	})
}

// Followers route
func Followers(c *gin.Context) {
	loggedIn(c, "")

	user := c.Param("id")
	username := CO.Get(user, "username")
	db := CO.DB()
	var followBy int
	followers := []interface{}{}
	me := CO.MeOrNot(c, user)
	var noMssg string

	stmt, _ := db.Prepare("SELECT followBy FROM follow WHERE followTo=? ORDER BY followID DESC")
	rows, fErr := stmt.Query(user)
	CO.Err(fErr)

	for rows.Next() {
		rows.Scan(&followBy)
		f := map[string]interface{}{
			"followBy": followBy,
		}
		followers = append(followers, f)
	}

	if me == true {
		noMssg = "You"
	} else {
		noMssg = username
	}

	renderTemplate(c, "followers", gin.H{
		"title":     username + "'s Followers",
		"session":   ses(c),
		"followers": followers,
		"no_mssg":   noMssg + " have no followers!!",
		"GET":       CO.Get,
		"UD":        CO.UsernameDecider,
		"noF":       CO.NoOfFollowers,
	})
}

// Followings route
func Followings(c *gin.Context) {
	loggedIn(c, "")

	user := c.Param("id")
	username := CO.Get(user, "username")
	db := CO.DB()
	var followTo int
	followings := []interface{}{}
	me := CO.MeOrNot(c, user)
	var noMssg string

	stmt, _ := db.Prepare("SELECT followTo FROM follow WHERE followBy=? ORDER BY followID DESC")
	rows, fErr := stmt.Query(user)
	CO.Err(fErr)

	for rows.Next() {
		rows.Scan(&followTo)
		f := map[string]interface{}{
			"followTo": followTo,
		}
		followings = append(followings, f)
	}

	if me == true {
		noMssg = "You"
	} else {
		noMssg = username
	}

	renderTemplate(c, "followings", gin.H{
		"title":      username + "'s Followings",
		"session":    ses(c),
		"followings": followings,
		"no_mssg":    noMssg + " have no followings!!",
		"GET":        CO.Get,
		"UD":         CO.UsernameDecider,
		"noF":        CO.NoOfFollowers,
	})
}

// Likes route
func Likes(c *gin.Context) {
	loggedIn(c, "")

	post := c.Param("id")
	db := CO.DB()
	var postCount int
	var likeBy int
	likes := []interface{}{}

	db.QueryRow("SELECT COUNT(postID) AS postCount FROM posts WHERE postID=?", post).Scan(&postCount)
	invalid(c, postCount)

	stmt, _ := db.Prepare("SELECT likeBy FROM likes WHERE postID=?")
	rows, err := stmt.Query(post)
	CO.Err(err)

	for rows.Next() {
		rows.Scan(&likeBy)
		l := map[string]interface{}{
			"likeBy": likeBy,
		}
		likes = append(likes, l)
	}

	renderTemplate(c, "likes", gin.H{
		"title":   "Likes",
		"session": ses(c),
		"likes":   likes,
		"GET":     CO.Get,
		"UD":      CO.UsernameDecider,
		"noF":     CO.NoOfFollowers,
	})
}
