package main

import (
	"example/gin-pro/handlers"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", handlers.Loginpage)
	r.GET("/register", handlers.Reigisterationpage)
	r.GET("/home", handlers.Homepage)
	r.POST("/login", handlers.Login)
	r.POST("/registration", handlers.Register)
	r.GET("/getproduct", handlers.Getproduct)
	r.POST("/deleteproduct/:ProductId", handlers.DeleteProduct)
	r.GET("/cart", handlers.Cartpage)
	r.GET("/addproduct", handlers.Addproductpage)
	r.POST("/addproduct", handlers.Addproduct)
	r.GET("/about", handlers.Aboutpage)
	r.GET("/users", handlers.GetuserPage)
	r.GET("/Getuserdata", handlers.Getuserdata)
	r.Run(":8085")
}
