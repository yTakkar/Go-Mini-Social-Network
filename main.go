package main

import (
	R "Go-Mini-Social-Network/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

func init() {
	godotenv.Load()
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	user := router.Group("/user")
	{
		user.POST("/signup", R.UserSignup)
		user.POST("/login", R.UserLogin)
	}

	router.GET("/", R.Index)
	router.GET("/welcome", R.Welcome)
	router.GET("/explore", R.Explore)
	router.GET("/404", R.NotFound)
	router.GET("/signup", R.Signup)
	router.GET("/login", R.Login)
	router.GET("/logout", R.Logout)
	router.GET("/edit_profile", R.EditProfile)
	router.GET("/create_post", R.CreatePost)

	router.GET("/profile/:id", R.Profile)
	router.GET("/profile", R.NotFound)

	router.GET("/view_post/:id", R.ViewPost)
	router.GET("/view_post", R.NotFound)

	router.GET("/edit_post/:id", R.EditPost)
	router.GET("/edit_post", R.NotFound)

	router.GET("/followers/:id", R.Followers)
	router.GET("/followers", R.NotFound)

	router.GET("/followings/:id", R.Followings)
	router.GET("/followings", R.NotFound)

	router.GET("/likes/:id", R.Likes)
	router.GET("/likes", R.NotFound)

	api := router.Group("/api")
	{
		api.POST("/create_new_post", R.CreateNewPost)
		api.POST("/delete_post", R.DeletePost)
		api.POST("/update_post", R.UpdatePost)
		api.POST("/update_profile", R.UpdateProfile)
		api.POST("/change_avatar", R.ChangeAvatar)
		api.POST("/follow", R.Follow)
		api.POST("/unfollow", R.Unfollow)
		api.POST("/like", R.Like)
		api.POST("/unlike", R.Unlike)
	}

	server := negroni.Classic()
	server.UseHandler(router)

	server.Run(os.Getenv("PORT"))

}
