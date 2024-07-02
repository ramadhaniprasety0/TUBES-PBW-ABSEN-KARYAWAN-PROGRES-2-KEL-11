package Usercontroller

import (
	"dashboard/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

var db *gorm.DB
var templates *template.Template

func InitUserController(database *gorm.DB, tmpl *template.Template) {
	db = database
	templates = tmpl
}

func Index(c *gin.Context) {
	var users []struct {
		ID   uint   `json:"id"`
		Nama string `json:"nama_admin"`
	}
	db.Model(&models.User{}).Select("id, nama").Scan(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Fungsi untuk hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Fungsi untuk verifikasi password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterUser(c *gin.Context) {
	nama := c.PostForm("nama")
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Hash password sebelum menyimpannya ke database
	hashedPassword, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Buat user baru
	user := models.User{Nama: nama, Username: username, Password: hashedPassword}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "redirect": "/"})
}

func LoginUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password Salah"})
		return
	}

	// Verifikasi password
	if !checkPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password Salah"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login Berhasil", "redirect": "/dashboard"})
}

func LogoutUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Save()
	c.Redirect(http.StatusMovedPermanently, "/")
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	c.Next()
}

func ShowDashboard(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	data := gin.H{
		"nama": user.Nama,
	}

	if err := templates.ExecuteTemplate(c.Writer, "index.html", data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

}
