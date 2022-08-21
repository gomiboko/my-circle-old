package main

import (
	"log"

	"github.com/gomiboko/my-circle/aws"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	err := db.Init()
	if err != nil {
		panic("failed to connect mycircle database!!")
	}

	err = aws.Init()
	if err != nil {
		panic("failed to load AWS SDK config!!")
	}

	r, err := server.NewRouter()
	if err != nil {
		panic("failed to generate router!!")
	}

	r.Run()
}
