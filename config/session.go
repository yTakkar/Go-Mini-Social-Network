package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

// GetSession to return the session
func GetSession(c *gin.Context) *sessions.Session {
	session, err := store.Get(c.Request, "session")
	Err(err)
	return session
}

// AllSessions function to return all the sessions
func AllSessions(c *gin.Context) (interface{}, interface{}) {
	session := GetSession(c)
	id := session.Values["id"]
	username := session.Values["username"]
	return id, username
}
