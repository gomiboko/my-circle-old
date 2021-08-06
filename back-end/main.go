package main

import (
	"log"

	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	err := db.Init()
	if err != nil {
		panic("failed to connect mycircle database!!")
	}

	r := server.NewRouter()
	r.Run()
}
