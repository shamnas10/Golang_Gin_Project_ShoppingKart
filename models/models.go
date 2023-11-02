package models

import "mime/multipart"

type User struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}
type Product struct {
	ProductName        string                `form:"productName" binding:"required"`
	ProductPrice       string                `form:"price" binding:"required"`
	ProductDescription string                `form:"productDescription" binding:"required"`
	ProductImage       *multipart.FileHeader `form:"image" binding:"required"`
}
type Getproduct struct {
	ProductId          int
	ProductName        string
	ProductPrice       string
	ProductDescription string
	ProductImage       string
}
type Userlist struct {
	Id       int
	Username string
	Email    string
}
