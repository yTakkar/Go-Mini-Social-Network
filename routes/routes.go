package routes

import (
	CO "Go-Mini-Social-Network/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hash(password string) []byte {
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	CO.Err(hashErr)
	return hash
}

func renderTemplate(c *gin.Context, tmpl string, p interface{}) {
	c.HTML(http.StatusOK, tmpl+".html", p)
}

func json(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func ses(c *gin.Context) interface{} {
	id, username := CO.AllSessions(c)
	return map[string]interface{}{
		"id":       id,
		"username": username,
	}
}

func loggedIn(c *gin.Context, urlRedirect string) {
	var URL string
	if urlRedirect == "" {
		URL = "/login"
	} else {
		URL = urlRedirect
	}
	id, _ := CO.AllSessions(c)
	if id == nil {
		c.Redirect(http.StatusFound, URL)
	}
}

func notLoggedIn(c *gin.Context) {
	id, _ := CO.AllSessions(c)
	if id != nil {
		c.Redirect(http.StatusFound, "/")
	}
}

func invalid(c *gin.Context, what int) {
	if what == 0 {
		c.Redirect(http.StatusFound, "/404")
	}
}
