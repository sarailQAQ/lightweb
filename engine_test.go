package lightweb

import (
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	r := Default()
	v := r.Group("/sigin")
	v.Use(mid)
	v.GET("/login",login)
	_,_,ok := r.router.getRouter("GET","/sigin/login")
	fmt.Println(ok)
	r.Run(8080)
}

func room(c *Context) {
	room := c.Param("room")
	fmt.Println("welcome to",room)
}

func TestParams(t *testing.T) {
	r := Default()
	r.GET("/sigin/:room",room)
	r.Run(8080)
}

func TestPost(t *testing.T) {
	r := Default()
	r.POST("/sigin", func(c *Context) {
		username := c.PostForm("username")
		passwd := c.PostForm("password")
		fmt.Println(username,passwd)
		c.JSON(200,H{"message":"hello"})
	})
	r.Run(8080)
}

