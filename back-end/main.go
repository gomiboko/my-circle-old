package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func main() {
	dbSvcNm := os.Getenv("DB_SERVICE_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbTz := os.Getenv("DB_TIME_ZONE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb3&parseTime=True&loc=%s", dbUser, dbPass, dbSvcNm, dbName, url.QueryEscape(dbTz))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect mycircle database!!")
	}

	var user User
	db.First(&user)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id":            user.ID,
			"name":          user.Name,
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
		})
	})
	r.Run()
}
