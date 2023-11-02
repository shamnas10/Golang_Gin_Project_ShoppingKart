package handlers

import (
	"database/sql"
	"example/gin-pro/database"
	"example/gin-pro/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Loginpage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Reigisterationpage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "homepage.html", nil)
}
func Login(c *gin.Context) {
	var loginData struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	var storedPassword string
	Db, err := database.Dbconnection()
	err = Db.QueryRow("SELECT password FROM users WHERE username = ?", loginData.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusOK, "index.html", gin.H{"Error": "Invalid Username"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if loginData.Password != storedPassword {
		c.HTML(http.StatusOK, "index.html", gin.H{"Error": "Invalid Password"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/home")
}

func Getuserdata(c *gin.Context) {
	pagen := c.DefaultQuery("page", "1")
	Page, _ := strconv.Atoi(pagen)
	pageSize := c.DefaultQuery("pageSize", "30")
	var Listusers []models.Userlist
	Db, err := database.UserDb()

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}
	offset := (Page - 1) * pageSizeNum
	// Use placeholders for limit and offset values in the query
	query := "SELECT * FROM newtable LIMIT ? OFFSET ?"
	rows, err := Db.Query(query, pageSizeNum, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Userlist
		if err := rows.Scan(&item.Id, &item.Username, &item.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		Listusers = append(Listusers, item)
	}
	c.JSON(http.StatusOK, gin.H{"Listusers": Listusers})
}

func GetuserPage(c *gin.Context) {

	c.HTML(http.StatusOK, "user.html", nil)
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Db, err := database.Dbconnection()

	// Insert user data into the database (Note: Password hashing is recommended)
	_, err = Db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		log.Println("Error inserting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"Error": "Registered Successfully"})

	c.Redirect(http.StatusSeeOther, "/")
}
func Getproduct(c *gin.Context) {
	var items []models.Getproduct
	Db, err := database.Dbconnection()
	rows, err := Db.Query("SELECT * FROM product")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Getproduct
		if err := rows.Scan(&item.ProductId, &item.ProductName, &item.ProductPrice, &item.ProductDescription, &item.ProductImage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	c.HTML(http.StatusOK, "myproduct.html", gin.H{
		"items": items, // Pass the items to the HTML template.
	})
}

func DeleteProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("ProductId"))
	if err != nil {
		// Log the value of ProductId for debugging
		fmt.Println("ProductId:", c.Param("ProductId"))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	Db, err := database.Dbconnection()

	// Execute a SQL DELETE statement to remove the product
	_, err = Db.Exec("DELETE FROM product WHERE id=?", productID)
	if err != nil {
		// Log the error for debugging
		fmt.Println("DELETE error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If the deletion was successful,
	c.Redirect(http.StatusSeeOther, "/getproduct")
}

func Cartpage(c *gin.Context) {
	c.HTML(http.StatusOK, "cart.html", nil)
}

func Addproductpage(c *gin.Context) {
	c.HTML(http.StatusOK, "addproduct.html", nil)
}
func Addproduct(c *gin.Context) {
	var product models.Product
	imageDirectory := "/home/shamnas/gin-pro/static/images/"

	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedFile, err := product.ProductImage.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer uploadedFile.Close()
	generatedUUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error generating UUID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}
	generatedFileName := generatedUUID.String() + filepath.Ext(product.ProductImage.Filename)

	// Define the full path, including the file name, for the uploaded image.
	imagePath := filepath.Join(imageDirectory, generatedFileName)

	// Create a file on the server to save the uploaded image.
	outputFile, err := os.Create(imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image file"})
		return
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, uploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image file"})
		return
	}
	imageRelativePath := strings.TrimPrefix(imagePath, "/home/shamnas/gin-pro")
	Db, err := database.Dbconnection()

	_, err = Db.Exec("INSERT INTO product (name, price, description, image) VALUES (?, ?, ?, ?)", product.ProductName, product.ProductPrice, product.ProductDescription, imageRelativePath)
	if err != nil {
		log.Println("Error inserting Product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to the database"})
		return
	}

	c.HTML(http.StatusOK, "addproduct.html", gin.H{"Error": "Product Added Successfully"})

}
func Aboutpage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}
