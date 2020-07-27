package lightweb

import (
	"fmt"
	"testing"
)

func login(ctx * Context) {
	username := ctx.Get("username")
	password := ctx.Query("password")
	fmt.Println(username,password)
	ctx.JSON(200,H{"message" : "Log in successful."})
}

func mid(ctx *Context) {
	username := ctx.Query("username")
	if username == "sarail" {
		fmt.Println("the force is with u")
	} else {
		fmt.Println("hello",username)
	}
	ctx.Set("username",username)
	ctx.Next()
}

func midd(ctx *Context) {
	username := ctx.Query("username")
	if username == "sarail" {
		fmt.Println("the force is with u")
	} else {
		fmt.Println("hello",username)
	}
	ctx.Set("username",username)
	ctx.JSON(500,H{"message":"Albort"})
	ctx.Abort()
}

func TestGet(t *testing.T) {
	r  := Default()
	r.GET("/login",mid,midd,login)
	r.Run(8080)
}


