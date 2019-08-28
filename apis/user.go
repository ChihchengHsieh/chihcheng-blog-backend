package apis

import (
	"blog/models"
	"blog/utils"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignupUser - Singup the user by "email", "password" and  "code"
func SignupUser(c *gin.Context) {
	// Extract required fields, including "email", "password" and "code"
	email, password, code := c.PostForm("email"), c.PostForm("password"), c.PostForm("code")

	// Checking if the password or password is empty
	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "You must provide email and password",
			"msg":      "You must provide email and password",
			"email":    email,
			"password": password,
		})
		return
	}

	// Checking if the register code is correct
	if code != os.Getenv("REGISTER_CODE") && code != os.Getenv("ADMIN_REGISTER_CODE") {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":  "The Given Rigister Code is not correct",
			"msg":  "The Given Rigister Code is not correct",
			"code": code,
		})
		return
	}

	// Checking if the user already exist
	if user, err := models.FindUserByEmail(email); user != nil {
		c.JSON(http.StatusConflict, gin.H{
			"err":  err,
			"msg":  "The user already exists.",
			"user": user,
		})
		return
	}

	// Using the register code to define the role
	role := "normal"
	if code == os.Getenv("ADMIN_REGISTER_CODE") {
		role = "admin"
	}

	// Hash the given password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err,
			"msg": "Cannot hash the given password",
		})
		return
	}

	// Create the newUser
	newUser := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	// Add this User
	insertedID, err := models.AddUser(&newUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"msg":     "Cannot register this user",
			"newUser": newUser,
		})
		return
	}

	newUser.ID = insertedID.(primitive.ObjectID)

	newUser.Password = ""

	authToken, err := utils.GenerateAuthToken(newUser.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"msg":   "Cannot generate the auth token for this user",
			"user":  newUser,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  newUser,
		"token": authToken,
	})

}

// LoginUser - Login the user through email and password
func LoginUser(c *gin.Context) {
	inputEmail, inputPassword := c.PostForm("email"), c.PostForm("password")

	if inputEmail == "" || inputPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Must Provide Email and Password",
			"msg":        "Must Provide Email and Password",
			"inputEmail": inputEmail,
		})
		return
	}

	user, err := models.CheckingTheAuth(inputEmail, inputPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      err,
			"msg":        "Email or Password is not correct",
			"user":       user,
			"inputEmail": inputEmail,
		})
		return
	}

	authToken, err := utils.GenerateAuthToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      err,
			"msg":        "Cannot Generate the Auth Token",
			"inputEmail": inputEmail,
			"user":       user,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": authToken,
		"user":  user,
	})

}
