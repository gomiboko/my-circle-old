package main

import (
	"net/http"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/server"
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
	log.SetFlags(log.LstdFlags | log.Llongfile)

	err := db.Init()
	if err != nil {
		panic("failed to connect mycircle database!!")
	}

	var user User
	db.GetDB().First(&user)

	r := server.NewRouter()
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
